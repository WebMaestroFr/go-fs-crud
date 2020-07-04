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

// parseFlags initializes flag variables
func parseFlags() {
	flag.StringVar(&port, "port", ":1234", "port to serve the API on")
	flag.StringVar(&store, "store", "./store", "path to the file storage directory")
	flag.Parse()
}

// getPath returns the path to a file
func getPath(name string) (path string) {
	path = filepath.Join(store, name)
	absPath, errPath := filepath.Abs(path)
	if errPath == nil {
		return absPath
	}
	return path
}

// getFilename returns the requested file name and path
func getFilename(request *http.Request) (name string, path string) {
	pathParams := mux.Vars(request)
	name = pathParams["name"]
	return name, getPath(name)
}

// getFilename returns the requested file name and path
func getFileResponse(name string, content []byte, path string) (f FileResponse) {
	fileInfo, _ := os.Stat(path)
	return FileResponse{name, string(content), path, fileInfo.Size(), fileInfo.ModTime()}
}

// handleError responds with a generic HTTP error
func handleError(w http.ResponseWriter, err error, code int) {
	errString := err.Error()
	http.Error(w, errString, code)
}

// handleResponse responds with a JSON file entry
func handleResponse(w http.ResponseWriter, response FileResponse) {
	r, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.Write(r)
}

// readFile returns the requested file instance or read operation error
func readFile(request *http.Request) (file FileResponse, errRead error) {
	var f FileResponse
	name, path := getFilename(request)
	content, errRead := ioutil.ReadFile(path)
	if errRead == nil {
		f = getFileResponse(name, content, path)
	}
	return f, errRead
}

// writeFile stores and returns the submitted file instance or write operation error
func writeFile(request *http.Request) (file FileResponse, errWrite error) {
	var f FileResponse
	name, path := getFilename(request)
	content, errContent := ioutil.ReadAll(request.Body)
	if errContent == nil {
		errWrite := ioutil.WriteFile(path, content, 0644)
		if errWrite == nil {
			f = getFileResponse(name, content, path)
		}
		return f, errWrite
	}
	return f, errContent
}

// removeFile removes the requested file or returns a remove operation error
func removeFile(request *http.Request) (errRemove error) {
	_, path := getFilename(request)
	return os.Remove(path)
}

// handleWrite is a generic "Write File" handler (for handleCreate and handleUpdate)
func handleWrite(w http.ResponseWriter, request *http.Request) {
	f, errWrite := writeFile(request)
	if errWrite == nil {
		handleResponse(w, f)
	} else {
		handleError(w, errWrite, 400)
	}
}

// handleCreate is the "Create File" handler for POST requests
func handleCreate(w http.ResponseWriter, request *http.Request) {
	_, errRead := readFile(request)
	if errRead == nil {
		handleError(w, errExists, 400)
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
		handleError(w, errRead, 400)
	}
}

// handleUpdate is the "Update File" handler for PUT requests
func handleUpdate(w http.ResponseWriter, request *http.Request) {
	_, errRead := readFile(request)
	if errRead == nil {
		handleWrite(w, request)
	} else {
		handleError(w, errRead, 400)
	}
}

// handleDelete is the "Delete File" handler for DELETE requests
func handleDelete(w http.ResponseWriter, request *http.Request) {
	errRemove := removeFile(request)
	if errRemove == nil {
		w.Write(msgRemove)
	} else {
		handleError(w, errRemove, 400)
	}
}

// main is the program main running process
func main() {
	// Initialize flag variables
	parseFlags()
	// Create a Router instance
	router := mux.NewRouter()
	// Define CRUD handlers
	router.HandleFunc("/{name}", handleCreate).Methods(http.MethodPost)
	router.HandleFunc("/{name}", handleRead).Methods(http.MethodGet)
	router.HandleFunc("/{name}", handleUpdate).Methods(http.MethodPut)
	router.HandleFunc("/{name}", handleDelete).Methods(http.MethodDelete)
	// Serve API
	server := http.ListenAndServe(port, router)
	// Log server errors to console
	log.Fatal(server)
}
