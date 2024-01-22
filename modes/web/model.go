package web

import (
	"net/http"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/amod"
	"github.com/asmaloney/gactar/framework"
)

var currentModelID = 1

type Model struct {
	id        int
	actrModel *actr.Model
}

type runOptions struct {
	Frameworks       []string `json:"frameworks,omitempty"` // list of frameworks to run on (if empty, "all")
	LogLevel         *string  `json:"logLevel,omitempty"`
	TraceActivations *bool    `json:"traceActivations,omitempty"`
	RandomSeed       *uint32  `json:"randomSeed,omitempty"`
}

func initModels(w *Web) {
	http.HandleFunc("/api/model/load", w.loadModelHandler)
}

func (w *Web) loadModelHandler(rw http.ResponseWriter, req *http.Request) {
	type request struct {
		SessionID int    `json:"sessionID"`
		AMODFile  string `json:"amod"`
	}
	type response struct {
		ModelID   int    `json:"modelID"`
		ModelName string `json:"modelName"`
		SessionID int    `json:"sessionID"`
	}

	var data request
	err := decodeBody(req, &data)
	if err != nil {
		encodeErrorResponse(rw, err)
		return
	}

	model, err := w.loadModel(data.SessionID, data.AMODFile)
	if err != nil {
		encodeErrorResponse(rw, err)
		return
	}

	encodeResponse(rw, response{
		ModelID:   model.id,
		ModelName: model.actrModel.Name,
		SessionID: data.SessionID,
	})
}

func (w *Web) loadModel(sessionID int, amodFile string) (model *Model, err error) {
	session := w.lookupSession(sessionID)
	if session == nil {
		err = &ErrInvalidSessionID{ID: sessionID}
		return
	}

	actrModel, err := generateModel(amodFile)
	if err != nil {
		return
	}

	model = &Model{
		id:        currentModelID,
		actrModel: actrModel,
	}
	currentModelID++

	session.addModel(model)

	return
}

// actrOptions converts runOptions into actr.Options
func actrOptions(options *runOptions) *actr.Options {
	if options == nil {
		return nil
	}

	// sets some defaults
	opts := actr.Options{
		LogLevel:         actr.ACTRLogLevel("info"),
		TraceActivations: false,
		RandomSeed:       options.RandomSeed,
	}

	if options.LogLevel != nil {
		opts.LogLevel = actr.ACTRLogLevel(*options.LogLevel)
	}

	if options.TraceActivations != nil {
		opts.TraceActivations = *options.TraceActivations
	}

	return &opts
}

func generateModel(amodFile string) (model *actr.Model, err error) {
	model, log, err := amod.GenerateModel(amodFile)
	if err != nil {
		err = &framework.ErrModelGenerationFailed{Log: log}
		return
	}

	return
}
