package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/gorilla/mux"
)

var router *mux.Router

func testRequest(t *testing.T, method string, url string, body io.Reader) *httptest.ResponseRecorder {
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rec := httptest.NewRecorder()
	// Create a request to pass to our handler. We don't have any query parameters for now, so we'll
	// pass 'nil' as the third parameter.
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		t.Fatal(err)
	} else {
		// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
		// directly and pass in our Request and ResponseRecorder.
		router.ServeHTTP(rec, req)
	}
	return rec
}

func TestCreate(t *testing.T) {
	body := strings.NewReader("Booyaka")
	res := testRequest(t, "POST", "/test.txt", body)
	// Check the status code is what we expect.
	if status := res.Code; status != http.StatusOK {
		t.Log(res.Body)
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestRead(t *testing.T) {
	res := testRequest(t, "GET", "/test.txt", nil)
	// Check the status code is what we expect.
	if status := res.Code; status != http.StatusOK {
		t.Log(res.Body)
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestUpdate(t *testing.T) {
	body := strings.NewReader("Boomshakalakasha")
	res := testRequest(t, "PUT", "/test.txt", body)
	// Check the status code is what we expect.
	if status := res.Code; status != http.StatusOK {
		t.Log(res.Body)
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestDelete(t *testing.T) {
	res := testRequest(t, "DELETE", "/test.txt", nil)
	// Check the status code is what we expect.
	if status := res.Code; status != http.StatusOK {
		t.Log(res.Body)
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestMain(m *testing.M) {
	// Initialize flag variables
	router = initializeRouter()
	// Run tests
	code := m.Run()
	os.Exit(code)
}
