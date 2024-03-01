package web

import (
	"net/http"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/amod"
	"github.com/asmaloney/gactar/framework"

	"github.com/asmaloney/gactar/util/runoptions"
)

var currentModelID = 1

type Model struct {
	id        int
	actrModel *actr.Model
}

// runOptionsJSON is the JSON version of runoptions.Options
type runOptionsJSON struct {
	Frameworks       runoptions.FrameworkNameList `json:"frameworks,omitempty"` // list of frameworks to run on (if empty, "all")
	LogLevel         *string                      `json:"logLevel,omitempty"`
	TraceActivations *bool                        `json:"traceActivations,omitempty"`
	RandomSeed       *uint32                      `json:"randomSeed,omitempty"`
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

// actrOptionsFromJSON converts runOptionsJSON into actr.Options. It defaults to the model's defaults.
func (w Web) actrOptionsFromJSON(defaults *runoptions.Options, options *runOptionsJSON) (*runoptions.Options, error) {
	if options == nil {
		return nil, nil
	}

	activeFrameworkNames := w.settings.Frameworks.Names()

	options.Frameworks.NormalizeFrameworkList(activeFrameworkNames)

	err := options.Frameworks.VerifyFrameworkList(activeFrameworkNames)
	if err != nil {
		return nil, err
	}

	opts := *defaults

	opts.Frameworks = options.Frameworks

	if options.LogLevel != nil {
		opts.LogLevel = runoptions.ACTRLogLevel(*options.LogLevel)
	}

	if options.TraceActivations != nil {
		opts.TraceActivations = *options.TraceActivations
	}

	opts.RandomSeed = options.RandomSeed

	return &opts, nil
}

func generateModel(amodFile string) (model *actr.Model, err error) {
	model, log, err := amod.GenerateModel(amodFile)
	if err != nil {
		err = &framework.ErrModelGenerationFailed{Log: log}
		return
	}

	return
}
