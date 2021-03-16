package sort

import (
	"bytes"
	"testing"
)

var inStr = "a\naa\nabbbb\naaaaaaaaaaaaaaaaaa\nabc"
var expStr = "a\naa\naaaaaaaaaaaaaaaaaa\nabbbb\nabc\n"

func TestLineSorter_Sort(t *testing.T) {
	in := bytes.NewBuffer([]byte(inStr))
	out := new(bytes.Buffer)
	sorter := lineSorter{
		r:        in,
		w:        out,
		buffSize: 10,
	}
	err := sorter.Sort()
	if err != nil {
		t.Fatal(err)
	}
	if expStr != out.String() {
		t.Fatal(out.String())
	}
}
