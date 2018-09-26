package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

func makeRouter() http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/health", func(http.ResponseWriter, *http.Request) {})

	return r
}
