package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
)

var store string

func storePath(filename string) string {
	return filepath.Join(store, filename)
}

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
	filename := storePath(name)
	_, errRead := ioutil.ReadFile(filename)
	if errRead == nil {
		errExists := errors.New("create: file exists already")
		handleError(w, errExists, 400)
	} else {
		content, errContent := ioutil.ReadAll(request.Body)
		if errContent == nil {
			errWrite := ioutil.WriteFile(filename, content, 0644)
			if errWrite == nil {
				handleRead(w, request)
			} else {
				handleError(w, errWrite, 400)
			}
		} else {
			handleError(w, errContent, 400)
		}
	}
}

func handleRead(w http.ResponseWriter, request *http.Request) {
	pathParams := mux.Vars(request)
	name := pathParams["name"]
	filename := storePath(name)
	content, errRead := ioutil.ReadFile(filename)
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
	filename := storePath(name)
	_, errRead := ioutil.ReadFile(filename)
	if errRead == nil {
		content, errContent := ioutil.ReadAll(request.Body)
		if errContent == nil {
			errWrite := ioutil.WriteFile(filename, content, 0644)
			if errWrite == nil {
				handleRead(w, request)
			} else {
				handleError(w, errWrite, 400)
			}
		} else {
			handleError(w, errContent, 400)
		}
	} else {
		handleError(w, errRead, 400)
	}
}

func handleDelete(w http.ResponseWriter, request *http.Request) {
	pathParams := mux.Vars(request)
	name := pathParams["name"]
	filename := storePath(name)
	errRemove := os.Remove(filename)
	if errRemove == nil {
		message := []byte("delete: successfully removed file")
		w.Write(message)
	} else {
		handleError(w, errRemove, 400)
	}
}

func main() {
	flag.StringVar(&store, "store", "./store", "path to the file storage directory")
	flag.Parse()
	router := mux.NewRouter()
	router.HandleFunc("/{name}", handleCreate).Methods(http.MethodPost)
	router.HandleFunc("/{name}", handleRead).Methods(http.MethodGet)
	router.HandleFunc("/{name}", handleUpdate).Methods(http.MethodPut)
	router.HandleFunc("/{name}", handleDelete).Methods(http.MethodDelete)
	server := http.ListenAndServe(":1234", router)
	log.Fatal(server)
}
