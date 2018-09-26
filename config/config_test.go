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
		t.Fatal(err)
	}
	defer os.Remove(file.Name())

	addr := "localhost:8080"
	timeout := 10
	url := "localhost:27017"
	config := fmt.Sprintf(`
[server]
address = "%s"
timeout = %d
[data_base]
url = "%s"
`, addr, timeout, url)

	_, err = file.Write([]byte(config))
	if err != nil {
		t.Fatal(err)
	}

	err = file.Close()
	if err != nil {
		t.Fatal(err)
	}

	conf := Make()
	err = conf.Open(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	ass.Equal(conf.Server.Address, addr)
	ass.Equal(conf.Server.Timeout, timeout)
	ass.Equal(conf.DataBase.URL, url)
}
