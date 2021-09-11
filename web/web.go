package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/urfave/cli/v2"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/amod"
	"github.com/asmaloney/gactar/framework"
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
}

func Initialize(cli *cli.Context, frameworks framework.List, examples *embed.FS) (w *Web, err error) {
	w = &Web{
		context:        cli,
		actrFrameworks: frameworks,
		examples:       examples,
		port:           cli.Int("port"),
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
		return nil, err
	}

	return
}

func (w *Web) Start() (err error) {
	http.HandleFunc("/run", w.runModel)

	exampleHandler := assetHandler(w.examples, "")
	http.HandleFunc("/examples/", exampleHandler.ServeHTTP)
	http.HandleFunc("/examples/list", w.listExamples)

	mainHandler := assetHandler(&mainAssets, "build")
	http.HandleFunc("/", mainHandler.ServeHTTP)

	fmt.Printf("Serving gactar on http://localhost:%d\n", w.port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", w.port), nil)
	if err != nil {
		return
	}

	return
}

func (w *Web) runModel(rw http.ResponseWriter, req *http.Request) {
	type request struct {
		AMODFile   string   `json:"amod"`
		RunStr     string   `json:"run"`
		Frameworks []string `json:"frameworks"`
	}
	type result struct {
		ModelName string `json:"modelName"`
		Code      string `json:"code"`
		Output    string `json:"output"`
	}

	type response struct {
		Results json.RawMessage `json:"results"`
	}

	decoder := json.NewDecoder(req.Body)

	var data request
	err := decoder.Decode(&data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	model, err := amod.GenerateModel(data.AMODFile)
	if err != nil {
		errorResponse(rw, err)
		return
	}

	resultMap := make(map[string]result, len(w.actrFrameworks))

	var wg sync.WaitGroup
	var mutex = &sync.Mutex{}

	for name, f := range w.actrFrameworks {
		wg.Add(1)

		go func(wg *sync.WaitGroup, name string, f framework.Framework) {
			defer wg.Done()

			code, output, err := w.run(model, data.RunStr, f)

			mutex.Lock()
			if err != nil {
				resultMap[name] = result{Output: err.Error()}
			} else {
				resultMap[name] = result{
					ModelName: model.Name,
					Code:      string(code),
					Output:    string(output),
				}
			}
			mutex.Unlock()

		}(&wg, name, f)
	}
	wg.Wait()

	results, err := json.Marshal(resultMap)
	if err != nil {
		errorResponse(rw, err)
		return
	}

	r := response{
		Results: json.RawMessage(string(results)),
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(r)
}

// assetHandler returns an http.Handler that will serve files from
// the given embed.FS.  When locating a file, it will prepend the root
// to the filesystem lookup.
// Adapted from https://blog.lawrencejones.dev/golang-embed/
func assetHandler(assets *embed.FS, root string) http.Handler {
	handler := fsFunc(func(name string) (fs.File, error) {
		assetPath := path.Join(root, name)

		f, err := assets.Open(assetPath)
		if os.IsNotExist(err) {
			return assets.Open("build/index.html")
		}

		return f, err
	})

	return http.FileServer(http.FS(handler))
}

// listExamples simply returns a list of the examples included in the build.
func (w *Web) listExamples(rw http.ResponseWriter, req *http.Request) {
	type response struct {
		List []string `json:"example_list"`
	}

	entries, err := w.examples.ReadDir("examples")
	if err != nil {
		errorResponse(rw, err)
		return
	}

	list := []string{}

	for _, entry := range entries {
		list = append(list, entry.Name())
	}

	r := response{
		List: list,
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(r)
}

func (w *Web) run(model *actr.Model, initialGoal string, framework framework.Framework) (generatedCode, output []byte, err error) {
	if model == nil {
		err = fmt.Errorf("no model loaded")
		return
	}

	err = framework.SetModel(model)
	if err != nil {
		return
	}

	initialGoal = strings.TrimSpace(initialGoal)

	generatedCode, output, err = framework.Run(initialGoal)
	if err != nil {
		return
	}

	return
}

func errorResponse(rw http.ResponseWriter, err error) {
	type response struct {
		ErrorStr string `json:"error"`
	}

	errResponse := response{
		ErrorStr: err.Error(),
	}

	json.NewEncoder(rw).Encode(errResponse)
}
