package web

import (
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitalyo61/genkeys/db"
	"github.com/vitalyo61/genkeys/generator"
)

func TestHandlers(t *testing.T) {
	ass := assert.New(t)

	db, err := db.Make("localhost:27017")
	ass.NoError(err)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	chData := make(chan *generator.Data)

	router := makeRouter(db, chData)

	gen, err := generator.Make("")
	ass.NoError(err)
	gen.Start(chData)
	defer func() {
		ch := make(chan *generator.Result)
		chData <- &generator.Data{
			Cmd:      generator.CmdStop,
			ChanCode: ch,
		}
		<-ch
	}()

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "/generate", nil)

	router.ServeHTTP(w, r)

	result := w.Result()
	body, _ := ioutil.ReadAll(result.Body)
	ass.Equal(result.StatusCode, http.StatusOK)
	ass.Equal(string(body), "0000")

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/generate", nil)

	router.ServeHTTP(w, r)

	result = w.Result()
	body, _ = ioutil.ReadAll(result.Body)
	ass.Equal(result.StatusCode, http.StatusOK)
	ass.Equal(string(body), "0001")

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/extinguish/0000", nil)

	router.ServeHTTP(w, r)

	result = w.Result()
	body, _ = ioutil.ReadAll(result.Body)
	ass.Equal(result.StatusCode, http.StatusOK)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/extinguish/0000", nil)

	router.ServeHTTP(w, r)

	result = w.Result()
	body, _ = ioutil.ReadAll(result.Body)
	ass.Equal(result.StatusCode, http.StatusInternalServerError)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/extinguish/0003", nil)

	router.ServeHTTP(w, r)

	result = w.Result()
	body, _ = ioutil.ReadAll(result.Body)
	ass.Equal(result.StatusCode, http.StatusInternalServerError)

	result = w.Result()
	body, _ = ioutil.ReadAll(result.Body)
	ass.Equal(result.StatusCode, http.StatusInternalServerError)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/info/0003", nil)

	router.ServeHTTP(w, r)

	result = w.Result()
	body, _ = ioutil.ReadAll(result.Body)
	ass.Equal(result.StatusCode, http.StatusOK)
	ass.Equal(string(body), "не выдан")

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/info/0001", nil)

	router.ServeHTTP(w, r)

	result = w.Result()
	body, _ = ioutil.ReadAll(result.Body)
	ass.Equal(result.StatusCode, http.StatusOK)
	ass.Equal(string(body), "выдан")

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/info/0000", nil)

	router.ServeHTTP(w, r)

	result = w.Result()
	body, _ = ioutil.ReadAll(result.Body)
	ass.Equal(result.StatusCode, http.StatusOK)
	ass.Equal(string(body), "погашен")

	count := uint32(math.Pow(62.0, 4.0) - 2)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/count", nil)

	router.ServeHTTP(w, r)

	result = w.Result()
	body, _ = ioutil.ReadAll(result.Body)
	ass.Equal(result.StatusCode, http.StatusOK)
	ass.Equal(string(body), fmt.Sprintf("%d", count))

	ch := make(chan *generator.Result)
	chData <- &generator.Data{
		Cmd:      generator.CmdSet,
		ChanCode: ch,
		Code:     "zzzz",
	}
	res := <-ch
	ass.NoError(res.Error)

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/count", nil)

	router.ServeHTTP(w, r)

	result = w.Result()
	body, _ = ioutil.ReadAll(result.Body)
	ass.Equal(result.StatusCode, http.StatusOK)
	ass.Equal(string(body), fmt.Sprintf("%d", 0))

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/generate", nil)

	router.ServeHTTP(w, r)

	result = w.Result()
	body, _ = ioutil.ReadAll(result.Body)
	ass.Equal(result.StatusCode, http.StatusServiceUnavailable)
	ass.Equal(string(body), generator.ErrorEndSequence.Error())

	err = db.CodeRemove("0000")
	ass.NoError(err)
	err = db.CodeRemove("0001")
	ass.NoError(err)
}
