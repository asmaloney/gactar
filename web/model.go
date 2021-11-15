package web

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/asmaloney/gactar/actr"
	"github.com/asmaloney/gactar/amod"
)

var currentModelID = 1

type Model struct {
	id        int
	actrModel *actr.Model
}

func initModels(w *Web) {
	http.HandleFunc("/model/load", w.loadModelHandler)
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
		err = fmt.Errorf("invalid session id '%d'", sessionID)
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

func generateModel(amodFile string) (model *actr.Model, err error) {
	model, log, err := amod.GenerateModel(amodFile)
	if err != nil {
		err = errors.New(log.String())
		return
	}

	return
}
