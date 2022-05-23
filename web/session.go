package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/asmaloney/gactar/framework"
)

type Session struct {
	id     int
	models []*Model
}

type SessionList []*Session

func initSessions(w *Web) {
	http.HandleFunc("/api/session/begin", w.beginSessionHandler)
	http.HandleFunc("/api/session/runModel", w.runModelSessionHandler)
	http.HandleFunc("/api/session/end", w.endSessionHandler)
}

func (w *Web) beginSessionHandler(rw http.ResponseWriter, req *http.Request) {
	type response struct {
		SessionID int `json:"session_id"`
	}

	session := w.newSession()

	encodeResponse(rw, response{
		SessionID: session.id,
	})
}

func (w *Web) runModelSessionHandler(rw http.ResponseWriter, req *http.Request) {
	type request struct {
		SessionID   int                      `json:"sessionID"`
		ModelID     int                      `json:"modelID"`
		Buffers     framework.InitialBuffers `json:"buffers"`              // set the initial buffers
		Frameworks  []string                 `json:"frameworks,omitempty"` // list of frameworks to run on (if empty, "all")
		IncludeCode bool                     `json:"includeCode"`          // include generated code in the result
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

	session := w.lookupSession(data.SessionID)
	if session == nil {
		err := fmt.Errorf("invalid session id '%d'", data.SessionID)
		encodeErrorResponse(rw, err)
		return
	}

	model := session.lookupModel(data.ModelID)
	if model == nil {
		err := fmt.Errorf("invalid model id '%d'", data.ModelID)
		encodeErrorResponse(rw, err)
		return
	}

	data.Frameworks = w.normalizeFrameworkList(data.Frameworks)

	err = w.verifyFrameworkList(data.Frameworks)
	if err != nil {
		encodeErrorResponse(rw, err)
		return
	}

	resultMap := w.runModel(model.actrModel, data.Buffers, data.Frameworks)

	for key := range resultMap {
		result := resultMap[key]

		// Remove the code if we just want the results
		if !data.IncludeCode {
			result.Code = nil
		}

		result.SessionID = &data.SessionID
		result.ModelID = &data.ModelID

		resultMap[key] = result
	}

	results, err := json.Marshal(resultMap)
	if err != nil {
		encodeErrorResponse(rw, err)
		return
	}

	encodeResponse(rw, response{
		Results: json.RawMessage(string(results)),
	})
}

func (w *Web) endSessionHandler(rw http.ResponseWriter, req *http.Request) {
	type request struct {
		SessionID int `json:"sessionID"`
	}
	type response struct {
	}

	var data request
	err := decodeBody(req, &data)
	if err != nil {
		encodeErrorResponse(rw, err)
		return
	}

	err = w.endSession(data.SessionID)
	if err != nil {
		encodeErrorResponse(rw, err)
		return
	}

	encodeResponse(rw, response{})
}

func (s *Session) addModel(model *Model) {
	s.models = append(s.models, model)
}

func (s *Session) lookupModel(modelID int) (model *Model) {
	for _, model := range s.models {
		if model.id == modelID {
			return model
		}
	}

	return nil
}

func (s *Session) end() {
	s.models = []*Model{}
}

func (w *Web) newSession() *Session {
	session := &Session{
		id: w.currentSessionID,
	}
	w.currentSessionID++

	w.sessionList = append(w.sessionList, session)

	return session
}

func (w *Web) endSession(id int) error {
	for index, session := range w.sessionList {
		if session.id == id {
			session.end()
			w.sessionList = removeSession(w.sessionList, index)
			return nil
		}
	}

	return fmt.Errorf("invalid session id '%d'", id)
}

func (w Web) hasSessions() bool {
	return len(w.sessionList) > 0
}

func (w Web) lookupSession(id int) *Session {
	for _, session := range w.sessionList {
		if session.id == id {
			return session
		}
	}

	return nil
}

func (w *Web) clearSessions() {
	for _, session := range w.sessionList {
		session.end()
	}

	w.sessionList = SessionList{}
}

func removeSession(s SessionList, index int) SessionList {
	return append(s[:index], s[index+1:]...)
}
