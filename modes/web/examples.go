package web

import "net/http"

func initExamples(w *Web) {
	exampleHandler := assetHandler(w.examples, "/api/examples/", "")
	http.HandleFunc("/api/examples/", exampleHandler.ServeHTTP)
	http.HandleFunc("/api/examples/list", w.listExamples)
}

// listExamples simply returns a list of the examples included in the build.
func (w *Web) listExamples(rw http.ResponseWriter, req *http.Request) {
	type response struct {
		List []string `json:"exampleList"`
	}

	entries, err := w.examples.ReadDir(".")
	if err != nil {
		encodeErrorResponse(rw, err)
		return
	}

	list := []string{}

	for _, entry := range entries {
		list = append(list, entry.Name())
	}

	encodeResponse(rw, response{
		List: list,
	})
}
