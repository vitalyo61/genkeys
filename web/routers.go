package web

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/vitalyo61/genkeys/db"
	"github.com/vitalyo61/genkeys/generator"
)

func makeRouter(b *db.DB, ch chan *generator.Data) http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/health", func(http.ResponseWriter, *http.Request) {})
	r.Use(setContentType)

	r.Handle("/generate", makeGenerateHandler(b, ch)).Methods(http.MethodGet)
	r.Handle("/extinguish/{code}", makeExtinguishHandler(b)).Methods(http.MethodGet)
	r.Handle("/info/{code}", makeInfoHandler(b)).Methods(http.MethodGet)
	r.Handle("/count", makeCountHandler(b, ch)).Methods(http.MethodGet)

	return r
}

func setContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}
