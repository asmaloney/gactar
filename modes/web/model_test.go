package web

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAddModel(t *testing.T) {
	session := webTest.newSession()

	if session == nil {
		t.Fatalf("Could not create session")
	}

	err := webTest.endSession(session.id)

	if err != nil {
		t.Errorf("Unexpected error: %s", err.Error())
	}

	if webTest.hasSessions() {
		t.Errorf("Did not remove session from list")
	}
}

func TestLoadModelHandler(t *testing.T) {
	session := webTest.newSession()

	src := `~~ model ~~
	name: Test
	~~ config ~~
	~~ init ~~
	~~ productions ~~`
	replacer := strings.NewReplacer(
		"\t", "",
		"\n", "\\n",
	)
	src = replacer.Replace(src)

	data := []byte(fmt.Sprintf(`{"sessionID":%d, "amod":"%s"}`, session.id, src))

	request, err := http.NewRequest("PUT", "/model/load", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(webTest.loadModelHandler)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned incorrect status code: expected '%v' got '%v'",
			http.StatusOK, status)
	}

	expected := `{"modelID":1,"modelName":"Test","sessionID":`
	responseStr := strings.TrimSpace(responseRecorder.Body.String())
	if !strings.HasPrefix(responseStr, expected) {
		t.Errorf("handler returned unexpected body: expected '%v' got '%v'",
			expected, responseStr)
	}

	if len(session.models) != 1 {
		t.Errorf("Model not loaded")
	}

	webTest.clearSessions()
}
