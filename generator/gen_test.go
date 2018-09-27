package generator

import (
	"math"
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
			current: "",
			next:    "0000",
			ok:      true,
		},
		{
			current: "0000",
			next:    "0001",
			ok:      true,
		},
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
		nextCode, ok := g.nextCode()
		ass.Equal(ok, tg.ok)
		ass.Equal(string(nextCode), tg.next)
	}
}

func TestCount(t *testing.T) {
	ass := assert.New(t)

	code := "zA90"
	count := countChar(code[3])
	ass.Equal(count, 0.0)
	count = countChar(code[2])
	ass.Equal(count, 9.0)
	count = countChar(code[1])
	ass.Equal(count, 10.0)
	count = countChar(code[0])
	ass.Equal(count, 61.0)

	c, err := Make("zzzz")
	ass.NoError(err)
	ass.Equal(c.freeCount(), uint32(0))

	err = c.set("zzzy")
	ass.NoError(err)
	ass.Equal(c.freeCount(), uint32(1))

	err = c.set("zzz0")
	ass.NoError(err)
	ass.Equal(c.freeCount(), uint32(math.Pow(62.0, 1.0)-1))

	err = c.set("zz00")
	ass.NoError(err)
	ass.Equal(c.freeCount(), uint32(math.Pow(62.0, 2.0)-1))

	err = c.set("0000")
	ass.NoError(err)
	ass.Equal(c.freeCount(), uint32(math.Pow(62.0, 4.0)-1))
}
