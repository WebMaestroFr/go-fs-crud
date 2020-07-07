// Package main starts the API server
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
	"time"

	"github.com/gorilla/mux"
)

// FileResponse type is a JSON format for file entries
type FileResponse struct {
	Name    string    `json:"name"`    // name and extension
	Content string    `json:"content"` // text content
	Path    string    `json:"path"`    // absolute path
	Size    int64     `json:"size"`    // length in bytes
	Time    time.Time `json:"time"`    // modification time
}

// Errors
var errExists = errors.New("create: file already exists")

// Messages
var msgRemove = []byte("delete: successfully removed file")

// Variables for flag values
var port string  // port to serve the API on
var store string // path to the file storage directory

// getPath returns the path to a file
func getPath(name string) (path string, err error) {
	path = filepath.Join(store, name)
	return filepath.Abs(path)
}

// getFilename returns the requested file name and path
func getFilename(request *http.Request) (name string, path string, err error) {
	pathParams := mux.Vars(request)
	name = pathParams["name"]
	path, errPath := getPath(name)
	return name, path, errPath
}

// getFilename returns the requested file name and path
func getFileResponse(name string, content []byte, path string) (f FileResponse, err error) {
	fileInfo, errStat := os.Stat(path)
	return FileResponse{name, string(content), path, fileInfo.Size(), fileInfo.ModTime()}, errStat
}

// handleError responds with a generic HTTP error
func handleError(w http.ResponseWriter, err error, code int) {
	errString := err.Error()
	http.Error(w, errString, code)
}

// handleResponse responds with a JSON file entry
func handleResponse(w http.ResponseWriter, response FileResponse) {
	r, errJSON := json.Marshal(response)
	if errJSON != nil {
		handleError(w, errJSON, http.StatusInternalServerError)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(r)
	}
}

// readFile returns the requested file instance or read operation error
func readFile(request *http.Request) (file FileResponse, err error) {
	var f FileResponse
	name, path, errFile := getFilename(request)
	if errFile != nil {
		return f, errFile
	}
	content, errRead := ioutil.ReadFile(path)
	if errRead != nil {
		return f, errRead
	}
	return getFileResponse(name, content, path)
}

// writeFile stores and returns the submitted file instance or write operation error
func writeFile(request *http.Request) (file FileResponse, err error) {
	var f FileResponse
	name, path, errFile := getFilename(request)
	if errFile != nil {
		return f, errFile
	}
	content, errContent := ioutil.ReadAll(request.Body)
	if errContent != nil {
		return f, errContent
	}
	errWrite := ioutil.WriteFile(path, content, 0644)
	if errWrite != nil {
		return f, errWrite
	}
	return getFileResponse(name, content, path)
}

// removeFile removes the requested file or returns a remove operation error
func removeFile(request *http.Request) (err error) {
	_, path, errFile := getFilename(request)
	if errFile != nil {
		return errFile
	}
	return os.Remove(path)
}

// handleWrite is a generic "Write File" handler (for handleCreate and handleUpdate)
func handleWrite(w http.ResponseWriter, request *http.Request) {
	f, errWrite := writeFile(request)
	if errWrite == nil {
		handleResponse(w, f)
	} else {
		handleError(w, errWrite, http.StatusBadRequest)
	}
}

// handleCreate is the "Create File" handler for POST requests
func handleCreate(w http.ResponseWriter, request *http.Request) {
	_, errRead := readFile(request)
	if errRead == nil {
		handleError(w, errExists, http.StatusBadRequest)
	} else {
		handleWrite(w, request)
	}
}

// handleRead is the "Read File" handler for GET requests
func handleRead(w http.ResponseWriter, request *http.Request) {
	f, errRead := readFile(request)
	if errRead == nil {
		handleResponse(w, f)
	} else {
		handleError(w, errRead, http.StatusBadRequest)
	}
}

// handleUpdate is the "Update File" handler for PUT requests
func handleUpdate(w http.ResponseWriter, request *http.Request) {
	_, errRead := readFile(request)
	if errRead == nil {
		handleWrite(w, request)
	} else {
		handleError(w, errRead, http.StatusBadRequest)
	}
}

// handleDelete is the "Delete File" handler for DELETE requests
func handleDelete(w http.ResponseWriter, request *http.Request) {
	errRemove := removeFile(request)
	if errRemove == nil {
		w.Write(msgRemove)
	} else {
		handleError(w, errRemove, http.StatusBadRequest)
	}
}

func initializeRouter() *mux.Router {
	// Initialize flag variables
	flag.StringVar(&port, "port", ":1234", "port to serve the API on")
	flag.StringVar(&store, "store", "/tmp/go-fs-crud", "path to the file storage directory")
	flag.Parse()
	// Initialize store directory
	errDir := os.Mkdir(store, 0744)
	if errDir != nil {
		log.Print(errDir)
	}
	// Create a Router instance
	router := mux.NewRouter()
	// Define CRUD handlers
	router.HandleFunc("/{name}", handleCreate).Methods(http.MethodPost)
	router.HandleFunc("/{name}", handleRead).Methods(http.MethodGet)
	router.HandleFunc("/{name}", handleUpdate).Methods(http.MethodPut)
	router.HandleFunc("/{name}", handleDelete).Methods(http.MethodDelete)
	return router
}

// main is the program main running process
func main() {
	// Create a Router instance
	router := initializeRouter()
	// Serve API
	server := http.ListenAndServe(port, router)
	// Log server errors to console
	log.Fatal(server)
}
