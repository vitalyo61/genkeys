package db

import (
	"testing"

	"github.com/globalsign/mgo"
	"github.com/stretchr/testify/assert"
)

func TestCode(t *testing.T) {
	ass := assert.New(t)

	db, err := Make("localhost:27017")
	ass.NoError(err)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	number := "test0000"
	code, err := db.CodeGet(number)
	ass.Error(mgo.ErrNotFound)

	code.Number = number
	err = db.CodeSet(code)
	ass.NoError(err)
	code, err = db.CodeGet(number)
	ass.NoError(err)
	ass.Equal(code.Number, number)
	ass.Equal(code.GetStatus(), "не выдан")

	err = db.CodeRemove(code.Number)
	ass.NoError(err)
}
