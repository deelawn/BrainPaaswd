package file

import (
	"io"
	"io/ioutil"
)

type Reader struct {
	source string
	done   bool
}

func (r *Reader) Read(p []byte) (n int, err error) {

	if r.done {
		return 0, io.EOF
	}

	data, err := ioutil.ReadFile(r.source)
	n = len(p)

	if err != nil {
		return n, err
	}

	for i, b := range data {
		p[i] = b
	}

	r.done = true
	return
}

func NewReader(source string) io.Reader {

	return &Reader{source: source}
}
