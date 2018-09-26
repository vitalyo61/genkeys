package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddChar(t *testing.T) {
	ass := assert.New(t)

	var next bool
	a := []byte("az")
	a[0], _ = addChar(a[0])
	ass.Equal([]byte("bz"), a)
	ass.False(next)

	a[1], next = addChar(a[1])
	ass.Equal([]byte("b0"), a)
	ass.True(next)
}

func TestGen(t *testing.T) {
	ass := assert.New(t)

	_, err := Make("qwert")
	ass.Error(err)

	_, err = Make("вася")
	ass.Error(err)

	testGen := []struct {
		current string
		next    string
		ok      bool
	}{
		{
			current: "zza9",
			next:    "zzaA",
			ok:      true,
		},
		{
			current: "zzaz",
			next:    "zzb0",
			ok:      true,
		},
		{
			current: "zzzz",
			next:    "0000",
			ok:      false,
		},
	}

	for _, tg := range testGen {
		g, err := Make(tg.current)
		if err != nil {
			t.Fatal(err.Error())
		}
		nextCode, ok := g.NextCode()
		ass.Equal(ok, tg.ok)
		ass.Equal(string(nextCode), tg.next)
	}
}
