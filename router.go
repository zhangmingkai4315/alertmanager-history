package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	// VERSION current build version
	VERSION = "0.1"
)

func versionHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Current Version %s", VERSION)
}

func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/version", versionHandler).Methods("GET")
	return r
}
