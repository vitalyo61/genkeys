package web

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vitalyo61/genkeys/config"
)

func TestServer(t *testing.T) {
	ass := assert.New(t)

	addr := ":8080"
	cfg := &config.Server{
		Address: addr,
		Timeout: 10,
	}
	srv := Make(cfg)
	go srv.Start()

	resp, err := http.Get("http://localhost:8080/health")
	if err != nil {
		t.Fatal(err)
	}
	ass.Equal(resp.StatusCode, http.StatusOK)

	err = srv.Stop()
	if err != nil {
		t.Fatal(err)
	}
}
