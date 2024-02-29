// Package web provides a web server with an HTTP API as well as a full UI to run amod code.
package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/jwalton/gchalk"
	"github.com/vearutop/statigz"
	"github.com/vearutop/statigz/brotli"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/amod"
	"github.com/asmaloney/gactar/framework"

	"github.com/asmaloney/gactar/util/cli"
	"github.com/asmaloney/gactar/util/container"
	"github.com/asmaloney/gactar/util/issues"
	"github.com/asmaloney/gactar/util/validate"
	"github.com/asmaloney/gactar/util/version"
)

//go:embed build/*
var mainAssets embed.FS

// fsFunc is used to access embedded files in assetHandler()
type fsFunc func(name string) (fs.File, error)

func (f fsFunc) Open(name string) (fs.File, error) {
	return f(name)
}

type Web struct {
	settings *cli.Settings
	examples *embed.FS
	port     int

	sessionList      SessionList
	currentSessionID int
}

type frameworkRunResult struct {
	ModelName string            `json:"modelName"`        // name of the model (from the amod file)
	Issues    *issues.IssueList `json:"issues,omitempty"` // issues specific to this framework

	FilePath *string `json:"filePath,omitempty"` // intermediate code file (full path)
	Code     *string `json:"code,omitempty"`     // actual code which was run
	Output   *string `json:"output,omitempty"`   // output of run (stdout + stderr)

	SessionID *int `json:"sessionID,omitempty"`
	ModelID   *int `json:"modelID,omitempty"`
}

type frameworkRunResultMap map[string]frameworkRunResult

type runResult struct {
	Issues  issues.IssueList      `json:"issues,omitempty"`
	Results frameworkRunResultMap `json:"results,omitempty"`
}

func Initialize(settings *cli.Settings, port int, examples *embed.FS) (w *Web, err error) {
	w = &Web{
		settings:         settings,
		examples:         examples,
		port:             port,
		sessionList:      SessionList{},
		currentSessionID: 1,
	}

	http.HandleFunc("/api/version", w.getVersionHandler)
	http.HandleFunc("/api/frameworks", w.getFrameworksHandler)
	http.HandleFunc("/api/run", w.runModelHandler)
	http.HandleFunc("/api/", http.NotFound)

	if examples != nil {
		initExamples(w)
	}

	initSessions(w)
	initModels(w)

	mainHandler := compressedAssetHandler(&mainAssets, "build")
	http.HandleFunc("/", mainHandler.ServeHTTP)

	return
}

func (w Web) Start() (err error) {
	fmt.Printf("Serving gactar on ")
	fmt.Println(gchalk.WithBlue().Underline(fmt.Sprintf("http://localhost:%d", w.port)))

	server := http.Server{
		Addr:        ":" + strconv.Itoa(w.port),
		ReadTimeout: 0,
	}

	err = server.ListenAndServe()
	if err != nil {
		return
	}

	return
}

func (Web) getVersionHandler(rw http.ResponseWriter, req *http.Request) {
	type response struct {
		Version string `json:"version"`
	}

	encodeResponse(rw, response{
		Version: version.BuildVersion,
	})
}

func (w Web) getFrameworksHandler(rw http.ResponseWriter, req *http.Request) {
	type response struct {
		Frameworks framework.InfoList `json:"frameworks"`
	}

	frameworks := framework.InfoList{}

	for _, framework := range w.settings.Frameworks {
		frameworks = append(frameworks, *framework.Info())
	}

	// return them sorted by name
	sort.Slice(frameworks, func(i, j int) bool {
		return frameworks[i].Name < frameworks[j].Name
	})

	encodeResponse(rw, response{
		Frameworks: frameworks,
	})
}

func (w Web) runModelHandler(rw http.ResponseWriter, req *http.Request) {
	type request struct {
		AMODFile string `json:"amod"` // text of an amod file
		Goal     string `json:"goal"` // initial goal

		Options *runOptionsJSON `json:"options,omitempty"`
	}

	var data request
	err := decodeBody(req, &data)
	if err != nil {
		encodeErrorResponse(rw, err)
		return
	}

	data.Options.Frameworks = w.normalizeFrameworkList(data.Options.Frameworks)

	err = w.verifyFrameworkList(data.Options.Frameworks)
	if err != nil {
		encodeErrorResponse(rw, err)
		return
	}

	model, log, err := amod.GenerateModel(data.AMODFile)
	if err != nil {
		encodeIssueResponse(rw, log)
		return
	}

	model.SetRunOptions(actrOptions(data.Options))

	initialGoal := strings.TrimSpace(data.Goal)
	initialBuffers := framework.InitialBuffers{
		"goal": initialGoal,
	}

	validate.Goal(model, initialGoal, log)

	// ensure temp dir exists
	// https://github.com/asmaloney/gactar/issues/103
	_, err = cli.CreateTempDir(w.settings)
	if err != nil {
		encodeErrorResponse(rw, err)
		return
	}

	resultMap := w.runModel(model, initialBuffers, data.Options.Frameworks)

	rr := runResult{
		Issues:  log.AllIssues(),
		Results: resultMap,
	}

	results, err := json.Marshal(rr)
	if err != nil {
		encodeErrorResponse(rw, err)
		return
	}

	encodeResponse(rw, json.RawMessage(string(results)))
}

// normalizeFrameworkList will look for "all" and replace it with all available
// framework names. It will then return a unique and sorted list of framework names.
func (w Web) normalizeFrameworkList(list []string) (normalized []string) {
	normalized = list

	if list == nil || slices.Contains(list, "all") {
		normalized = w.settings.Frameworks.Names()
	}

	normalized = container.UniqueAndSorted(normalized)
	return
}

// verifyFrameworkList will check that each name is of a valid framework and that
// it is active on this server.
func (w Web) verifyFrameworkList(list []string) (err error) {
	for _, name := range list {
		if !framework.IsValidFramework(name) {
			err = &ErrInvalidFrameworkName{Name: name}
			return
		}

		// we have a valid name, check if it is active
		if _, ok := w.settings.Frameworks[name]; !ok {
			err = &ErrFrameworkNotActive{Name: name}
			return
		}
	}

	return
}

func (w Web) runModel(model *actr.Model, initialBuffers framework.InitialBuffers, frameworkNames []string) (resultMap frameworkRunResultMap) {
	resultMap = make(frameworkRunResultMap, len(frameworkNames))

	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}

	for _, name := range frameworkNames {
		f := w.settings.Frameworks[name]

		wg.Add(1)

		go func(wg *sync.WaitGroup, name string, f framework.Framework) {
			defer wg.Done()

			result := &framework.RunResult{}

			log := f.ValidateModel(model)
			if !log.HasError() {
				r, err := runModelOnFramework(model, initialBuffers, f)
				if err != nil {
					log.Error(nil, err.Error())
				}
				if r != nil {
					result = r
				}
			}

			frameworkResult := frameworkRunResult{
				ModelName: model.Name,
			}

			mutex.Lock()

			if log.HasIssues() {
				all := log.AllIssues()
				frameworkResult.Issues = &all
			}

			if result.FileName != "" {
				frameworkResult.FilePath = &result.FileName
			}

			if len(result.GeneratedCode) > 0 {
				codeStr := string(result.GeneratedCode)
				frameworkResult.Code = &codeStr

			}
			if len(result.Output) > 0 {
				outputStr := string(result.Output)
				frameworkResult.Output = &outputStr

			}

			resultMap[name] = frameworkResult

			mutex.Unlock()
		}(&wg, name, f)
	}
	wg.Wait()

	return
}

func runModelOnFramework(model *actr.Model, initialBuffers framework.InitialBuffers, f framework.Framework) (result *framework.RunResult, err error) {
	if model == nil {
		err = ErrNoModel
		return
	}

	err = f.SetModel(model)
	if err != nil {
		return
	}

	result, err = f.Run(initialBuffers)
	if err != nil {
		return
	}

	return
}

func decodeBody(req *http.Request, v interface{}) (err error) {
	if req.Body == nil {
		return ErrEmptyRequestBody
	}

	decoder := json.NewDecoder(req.Body)

	err = decoder.Decode(&v)
	if err != nil {
		return err
	}

	return
}

func encodeResponse(rw http.ResponseWriter, v interface{}) {
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	encodeErr := json.NewEncoder(rw).Encode(v)
	if encodeErr != nil {
		http.Error(rw, encodeErr.Error(), http.StatusInternalServerError)
	}
}

func encodeErrorResponse(rw http.ResponseWriter, err error) {
	errResponse := runResult{
		Issues: issues.IssueList{
			{
				Level: "error",
				Text:  err.Error(),
			},
		},
	}

	encodeErr := json.NewEncoder(rw).Encode(errResponse)
	if encodeErr != nil {
		http.Error(rw, encodeErr.Error(), http.StatusInternalServerError)
	}
}

func encodeIssueResponse(rw http.ResponseWriter, log *issues.Log) {
	errResponse := runResult{Issues: log.AllIssues()}

	encodeErr := json.NewEncoder(rw).Encode(errResponse)
	if encodeErr != nil {
		http.Error(rw, encodeErr.Error(), http.StatusInternalServerError)
	}
}

// assetHandler returns an http.Handler that will serve files from
// the given embed.FS.  When locating a file, it will optionally strip
// and append a prefix to the filesystem lookup.
// Adapted from https://blog.lawrencejones.dev/golang-embed/
func assetHandler(assets *embed.FS, stripPrefix, prefix string) http.Handler {
	handler := fsFunc(func(name string) (fs.File, error) {
		assetPath := path.Join(prefix, name)

		f, err := assets.Open(assetPath)
		if os.IsNotExist(err) {
			return assets.Open("build/index.html")
		}

		return f, err
	})

	return http.StripPrefix(stripPrefix, http.FileServer(http.FS(handler)))
}

// compressedAssetHandler returns an http.Handler that will serve files from
// the given embed.FS using statigz to serve compressed files.
func compressedAssetHandler(assets *embed.FS, prefix string) http.Handler {
	return statigz.FileServer(assets, brotli.AddEncoding, statigz.FSPrefix(prefix))
}
