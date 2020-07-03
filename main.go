package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

type fileResponse struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

func handleError(w http.ResponseWriter, err error, code int) {
	http.Error(w, err.Error(), code)
}

func handleFileResponse(w http.ResponseWriter, response fileResponse) {
	r, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(r)
}

func handleCreate(w http.ResponseWriter, request *http.Request) {
	pathParams := mux.Vars(request)
	name := pathParams["name"]
	_, errRead := ioutil.ReadFile(name)
	if errRead == nil {
		errExists := errors.New("create: file exists already")
		handleError(w, errExists, 400)
	} else {
		errWrite := ioutil.WriteFile(name, []byte("Test"), 0644)
		if errWrite == nil {
			handleRead(w, request)
		} else {
			handleError(w, errWrite, 400)
		}
	}
}

func handleRead(w http.ResponseWriter, request *http.Request) {
	pathParams := mux.Vars(request)
	name := pathParams["name"]
	content, errRead := ioutil.ReadFile(name)
	if errRead == nil {
		response := fileResponse{name, string(content)}
		handleFileResponse(w, response)
	} else {
		handleError(w, errRead, 400)
	}
}

func handleUpdate(w http.ResponseWriter, request *http.Request) {
	pathParams := mux.Vars(request)
	name := pathParams["name"]
	_, errRead := ioutil.ReadFile(name)
	if errRead == nil {
		errWrite := ioutil.WriteFile(name, []byte("Test Updated"), 0644)
		if errWrite == nil {
			handleRead(w, request)
		} else {
			handleError(w, errWrite, 400)
		}
	} else {
		handleError(w, errRead, 400)
	}
}

func handleDelete(w http.ResponseWriter, request *http.Request) {
	pathParams := mux.Vars(request)
	name := pathParams["name"]
	errRemove := os.Remove(name)
	if errRemove == nil {
		message := []byte("delete: successfully removed file")
		w.Write(message)
	} else {
		handleError(w, errRemove, 400)
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/{name}", handleCreate).Methods(http.MethodPost)
	router.HandleFunc("/{name}", handleRead).Methods(http.MethodGet)
	router.HandleFunc("/{name}", handleUpdate).Methods(http.MethodPut)
	router.HandleFunc("/{name}", handleDelete).Methods(http.MethodDelete)
	server := http.ListenAndServe(":1234", router)
	log.Fatal(server)
}
