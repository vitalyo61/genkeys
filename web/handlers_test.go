package web

import (
	"io/ioutil"
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
	r := httptest.NewRequest(http.MethodGet, "/get", nil)

	router.ServeHTTP(w, r)

	result := w.Result()
	body, _ := ioutil.ReadAll(result.Body)
	ass.Equal(result.StatusCode, http.StatusOK)
	ass.Equal(string(body), "0000")

	w = httptest.NewRecorder()
	r = httptest.NewRequest(http.MethodGet, "/get", nil)

	router.ServeHTTP(w, r)

	result = w.Result()
	body, _ = ioutil.ReadAll(result.Body)
	ass.Equal(result.StatusCode, http.StatusOK)
	ass.Equal(string(body), "0001")

	err = db.CodeRemove("0000")
	ass.NoError(err)
	err = db.CodeRemove("0001")
	ass.NoError(err)
}
