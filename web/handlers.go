package web

import (
	"log"
	"net/http"

	"github.com/globalsign/mgo"
	"github.com/gorilla/mux"
	"github.com/vitalyo61/genkeys/db"
	"github.com/vitalyo61/genkeys/db/model"
	"github.com/vitalyo61/genkeys/generator"
)

type generateHandler struct {
	dBase   *db.DB
	chGen   chan *generator.Data
	errHead string
}

func makeGenerateHandler(b *db.DB, ch chan *generator.Data) http.Handler {
	return &generateHandler{
		dBase:   b,
		chGen:   ch,
		errHead: "generate code:",
	}
}

func (h *generateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	chResult := make(chan *generator.Result)
	data := &generator.Data{
		Cmd:      generator.CmdGet,
		ChanCode: chResult,
	}

	h.chGen <- data
	res := <-chResult

	if res.Error != nil {
		log.Printf("%s %s\n", h.errHead, res.Error)
		w.WriteHeader(http.StatusServiceUnavailable)
		_, err = w.Write([]byte(res.Error.Error()))
	} else {
		code := &model.Code{
			Number: string(res.Code),
			Status: model.CodeYes,
		}

		err = h.dBase.CodeSet(code)
		if err != nil {
			log.Printf("%s %s\n", h.errHead, err)
			w.WriteHeader(http.StatusInternalServerError)
			_, err = w.Write([]byte(err.Error()))
		} else {
			_, err = w.Write(res.Code)
		}
	}

	if err != nil {
		log.Printf("%s %s\n", h.errHead, err)
	}
}

type extinguishHandler struct {
	dBase   *db.DB
	errHead string
}

func makeExtinguishHandler(b *db.DB) http.Handler {
	return &extinguishHandler{
		dBase:   b,
		errHead: "extinguish code:",
	}
}

func (h *extinguishHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	number := mux.Vars(r)["code"]

	err := h.dBase.CodeExtinguish(number)
	if err != nil {
		log.Printf("%s %s\n", h.errHead, err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
	}

	if err != nil {
		log.Printf("%s %s\n", h.errHead, err)
	}
}

type infoHandler struct {
	dBase   *db.DB
	errHead string
}

func makeInfoHandler(b *db.DB) http.Handler {
	return &infoHandler{
		dBase:   b,
		errHead: "info code:",
	}
}

func (h *infoHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	number := mux.Vars(r)["code"]

	code, err := h.dBase.CodeGet(number)
	if err != nil && err != mgo.ErrNotFound {
		log.Printf("%s %s\n", h.errHead, err)
		w.WriteHeader(http.StatusInternalServerError)
		_, err = w.Write([]byte(err.Error()))
	} else {
		_, err = w.Write([]byte(code.GetStatus()))
	}

	if err != nil {
		log.Printf("%s %s\n", h.errHead, err)
	}
}

type countHandler struct {
	dBase   *db.DB
	chGen   chan *generator.Data
	errHead string
}

func makeCountHandler(b *db.DB, ch chan *generator.Data) http.Handler {
	return &countHandler{
		dBase:   b,
		chGen:   ch,
		errHead: "count code:",
	}
}

func (h *countHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var err error

	chResult := make(chan *generator.Result)
	data := &generator.Data{
		Cmd:      generator.CmdCount,
		ChanCode: chResult,
	}

	h.chGen <- data
	res := <-chResult

	_, err = w.Write(res.Code)

	if err != nil {
		log.Printf("%s %s\n", h.errHead, err)
	}
}
