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
	"sort"
	"strings"
	"sync"

	"github.com/urfave/cli/v2"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/amod"
	"github.com/asmaloney/gactar/framework"
	"github.com/asmaloney/gactar/version"

	"github.com/asmaloney/gactar/util/container"
	"github.com/asmaloney/gactar/util/filesystem"
	"github.com/asmaloney/gactar/util/issues"
	"github.com/asmaloney/gactar/util/validate"
)

//go:embed build/*
var mainAssets embed.FS

// fsFunc is used to access embedded files in assetHandler()
type fsFunc func(name string) (fs.File, error)

func (f fsFunc) Open(name string) (fs.File, error) {
	return f(name)
}

type Web struct {
	context        *cli.Context
	actrFrameworks framework.List
	examples       *embed.FS
	port           int

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

func Initialize(cli *cli.Context, frameworks framework.List, examples *embed.FS) (w *Web, err error) {
	w = &Web{
		context:          cli,
		actrFrameworks:   frameworks,
		examples:         examples,
		port:             cli.Int("port"),
		sessionList:      SessionList{},
		currentSessionID: 1,
	}

	for name, f := range w.actrFrameworks {
		err = f.Initialize()
		if err != nil {
			fmt.Println(err.Error())
			delete(w.actrFrameworks, name)
			err = nil
		}
	}

	if len(w.actrFrameworks) == 0 {
		err := fmt.Errorf("could not initialize any frameworks - please check your installation")
		return w, err
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

	mainHandler := assetHandler(&mainAssets, "", "build")
	http.HandleFunc("/", mainHandler.ServeHTTP)

	return
}

func (w Web) Start() (err error) {
	fmt.Printf("Serving gactar on http://localhost:%d\n", w.port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", w.port), nil)
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

	for _, framework := range w.actrFrameworks {
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
		AMODFile   string   `json:"amod"`                 // text of an amod file
		Goal       string   `json:"goal"`                 // initial goal
		Frameworks []string `json:"frameworks,omitempty"` // list of frameworks to run on (if empty, "all")
	}

	type response struct {
		Results json.RawMessage `json:"results"`
	}

	var data request
	err := decodeBody(req, &data)
	if err != nil {
		encodeErrorResponse(rw, err)
		return
	}

	data.Frameworks = w.normalizeFrameworkList(data.Frameworks)

	err = w.verifyFrameworkList(data.Frameworks)
	if err != nil {
		encodeErrorResponse(rw, err)
		return
	}

	model, log, err := amod.GenerateModel(data.AMODFile)
	if err != nil {
		encodeIssueResponse(rw, log)
		return
	}

	initialGoal := strings.TrimSpace(data.Goal)
	initialBuffers := framework.InitialBuffers{
		"goal": initialGoal,
	}

	validate.Goal(model, initialGoal, log)

	resultMap := w.runModel(model, initialBuffers, data.Frameworks)

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

	if list == nil || container.Contains("all", list) {
		normalized = w.actrFrameworks.Names()
	}

	normalized = container.UniqueAndSorted(normalized)
	return
}

// verifyFrameworkList will check that each name is of a valid framework and that
// it is active on this server.
func (w Web) verifyFrameworkList(list []string) (err error) {
	for _, name := range list {
		if !framework.IsValidFramework(name) {
			err = fmt.Errorf("invalid framework name: %q", name)
			return
		}

		// we have a valid name, check if it is active
		if _, ok := w.actrFrameworks[name]; !ok {
			err = fmt.Errorf("framework %q is not active on server", name)
			return
		}
	}

	return
}

func (w Web) runModel(model *actr.Model, initialBuffers framework.InitialBuffers, frameworkNames []string) (resultMap frameworkRunResultMap) {
	// ensure temp dir exists
	// https://github.com/asmaloney/gactar/issues/103
	filesystem.CreateTempDir(w.context)

	resultMap = make(frameworkRunResultMap, len(frameworkNames))

	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}

	for _, name := range frameworkNames {
		f := w.actrFrameworks[name]

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
		err = fmt.Errorf("no model loaded")
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
		err = fmt.Errorf("empty request body")
		return err
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
	json.NewEncoder(rw).Encode(v)
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

	json.NewEncoder(rw).Encode(errResponse)
}

func encodeIssueResponse(rw http.ResponseWriter, log *issues.Log) {
	errResponse := runResult{Issues: log.AllIssues()}

	json.NewEncoder(rw).Encode(errResponse)
}

// assetHandler returns an http.Handler that will serve files from
// the given embed.FS.  When locating a file, it will optionally strip
// and append a prefix to the filesystem lookup.
// Adapted from https://blog.lawrencejones.dev/golang-embed/
func assetHandler(assets *embed.FS, stripPrefix, prepend string) http.Handler {
	handler := fsFunc(func(name string) (fs.File, error) {
		assetPath := path.Join(prepend, name)

		f, err := assets.Open(assetPath)
		if os.IsNotExist(err) {
			return assets.Open("build/index.html")
		}

		return f, err
	})

	return http.StripPrefix(stripPrefix, http.FileServer(http.FS(handler)))
}
