package web

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"net/http"
	"os"
	"path"

	"github.com/urfave/cli/v2"

	"gitlab.com/asmaloney/gactar/actr"
	"gitlab.com/asmaloney/gactar/amod"
	"gitlab.com/asmaloney/gactar/framework"
)

//go:embed build/*
var assets embed.FS

// fsFunc is used to access embedded files in assetHandler()
type fsFunc func(name string) (fs.File, error)

func (f fsFunc) Open(name string) (fs.File, error) {
	return f(name)
}

type Web struct {
	context       *cli.Context
	actrFramework framework.Framework
	port          int
}

func Initialize(cli *cli.Context, framework framework.Framework) (w *Web, err error) {
	w = &Web{
		context:       cli,
		actrFramework: framework,
		port:          cli.Int("port"),
	}

	err = framework.Initialize()
	if err != nil {
		return nil, err
	}

	return
}

func (w *Web) Start() (err error) {
	http.HandleFunc("/run", w.runModel)

	handler := assetHandler("build")
	http.HandleFunc("/", handler.ServeHTTP)

	fmt.Printf("Serving gactar on http://localhost:%d\n", w.port)

	err = http.ListenAndServe(fmt.Sprintf(":%d", w.port), nil)
	if err != nil {
		return
	}

	return
}

func (w *Web) runModel(rw http.ResponseWriter, req *http.Request) {
	type request struct {
		AMODFile string `json:"amod"`
		RunStr   string `json:"run"`
	}

	type response struct {
		Results string `json:"results"`
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

	output, err := w.run(model, data.RunStr)
	if err != nil {
		errorResponse(rw, err)
		return
	}

	r := response{
		Results: string(output),
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	json.NewEncoder(rw).Encode(r)
}

// assetHandler returns an http.Handler that will serve files from
// the assets embed.FS.  When locating a file, it will prepend the root
// to the filesystem lookup.
// From https://blog.lawrencejones.dev/golang-embed/ and modified.
func assetHandler(root string) http.Handler {
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

func (w *Web) run(model *actr.Model, initialGoal string) (output []byte, err error) {
	if model == nil {
		err = fmt.Errorf("no model loaded")
		return
	}

	err = w.actrFramework.SetModel(model)
	if err != nil {
		return
	}

	output, err = w.actrFramework.Run(initialGoal)
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
