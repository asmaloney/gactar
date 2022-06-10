package web

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewSession(t *testing.T) {
	session := webTest.newSession()

	if session == nil {
		t.Errorf("Could not create session")
	}

	webTest.clearSessions()
}

func TestEndSession(t *testing.T) {
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

func TestBeginSessionHandler(t *testing.T) {
	request, err := http.NewRequest("PUT", "/session/begin", nil)
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(webTest.beginSessionHandler)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned incorrect status code: expected '%v' want '%v'",
			http.StatusOK, status)
	}

	expected := `{"session_id":`
	responseStr := strings.TrimSpace(responseRecorder.Body.String())
	if !strings.HasPrefix(responseStr, expected) {
		t.Errorf("handler returned unexpected body: expected to start with '%v' got '%v'",
			expected, responseStr)
	}

	webTest.clearSessions()
}

func TestEndSessionHandler(t *testing.T) {
	session := webTest.newSession()

	data := []byte(fmt.Sprintf(`{"sessionID":%d}`, session.id))

	request, err := http.NewRequest("PUT", "/session/end", bytes.NewBuffer(data))
	if err != nil {
		t.Fatal(err)
	}

	responseRecorder := httptest.NewRecorder()
	handler := http.HandlerFunc(webTest.endSessionHandler)

	handler.ServeHTTP(responseRecorder, request)

	if status := responseRecorder.Code; status != http.StatusOK {
		t.Errorf("handler returned incorrect status code: expected '%v' got '%v'",
			http.StatusOK, status)
	}

	expected := "{}"
	responseStr := strings.TrimSpace(responseRecorder.Body.String())
	if !strings.HasPrefix(responseStr, expected) {
		t.Errorf("handler returned unexpected body: expected '%v' got '%v'",
			expected, responseStr)
	}

	if webTest.hasSessions() {
		t.Errorf("Did not remove session from list")
	}
}

// Commented out for now since the CI does not install any frameworks.

// func TestRunModelSessionHandler(t *testing.T) {
// 	session := webTest.newSession()

// 	src := `==model==
// 	name: Test
// 	==config==
// 	gactar { log_level: 'min' }
// 	chunks {
// 		[count: first second]
// 		[countFrom: start end status]
// 	}
// 	==init==
// 	memory {
// 		[count: 0 1]
// 		[count: 1 2]
// 		[count: 2 3]
// 		[count: 3 4]
// 		[count: 4 5]
// 		[count: 5 6]
// 		[count: 6 7]
// 		[count: 7 8]
// 	}
// 	==productions==
// 	start {
// 		match {
// 			goal [countFrom: ?start ?end starting]
// 		}
// 		do {
// 			recall [count: ?start ?]
// 			set goal to [countFrom: ?start ?end counting]
// 		}
// 	}
// 	increment {
// 		match {
// 			goal [countFrom: ?x !?x counting]
// 			retrieval [count: ?x ?next]
// 		}
// 		do {
// 			print ?x
// 			recall [count: ?next ?]
// 			set goal.start to ?next
// 		}
// 	}
// 	stop {
// 		match {
// 			goal [countFrom: ?x ?x counting]
// 		}
// 		do {
// 			print ?x
// 			clear goal
// 		}
// 	}`

// 	model, err := webTest.loadModel(session.id, src)
// 	if err != nil {
// 		t.Errorf("Unexpected error: %s", err.Error())
// 		return
// 	}

// 	data := []byte(fmt.Sprintf(`{"sessionID":%d, "modelID":%d, "buffers":{ "goal":"[countFrom: 2 5 starting]" }}`, session.id, model.id))

// 	request, err := http.NewRequest("PUT", "/session/run", bytes.NewBuffer(data))
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	responseRecorder := httptest.NewRecorder()
// 	handler := http.HandlerFunc(webTest.runModelSessionHandler)

// 	handler.ServeHTTP(responseRecorder, request)

// 	if status := responseRecorder.Code; status != http.StatusOK {
// 		t.Errorf("handler returned incorrect status code: expected '%v' got '%v'",
// 			http.StatusOK, status)
// 	}

// 	expected := `{"results":{"ccm":{"language":"python","modelName":"Test","filePath":"ccm_Test.py","output":"2\n3\n4\n5\nend...\n"`
// 	responseStr := strings.TrimSpace(responseRecorder.Body.String())
// 	if !strings.HasPrefix(responseStr, expected) {
// 		t.Errorf("handler returned unexpected body: expected '%v' got '%v'",
// 			expected, responseStr)
// 	}

// 	webTest.endSession(session.id)

// 	if webTest.hasSessions() {
// 		t.Errorf("Did not remove session from list")
// 	}
// }
