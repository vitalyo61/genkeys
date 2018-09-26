package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConfig(t *testing.T) {
	ass := assert.New(t)

	file, err := ioutil.TempFile("", "temp")
	if err != nil {
		t.Fatal(err.Error())
	}
	defer os.Remove(file.Name())

	addr := "localhost:8080"
	timeout := 10
	config := fmt.Sprintf(`
[server]
address = "%s"
timeout = %d
`, addr, timeout)

	_, err = file.Write([]byte(config))
	if err != nil {
		t.Fatal(err.Error())
	}

	err = file.Close()
	if err != nil {
		t.Fatal(err.Error())
	}

	conf := Make()
	err = conf.Open(file.Name())
	if err != nil {
		t.Fatal(err.Error())
	}

	ass.Equal(conf.Server.Address, addr)
	ass.Equal(conf.Server.Timeout, timeout)
}
