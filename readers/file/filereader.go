package file

import (
	"io"
	"io/ioutil"
	"os"
	"time"

	"github.com/deelawn/BrainPaaswd/readers"
)

// Reader is a file reader that implements io.Reader
type Reader struct {
	source string
	done   bool
}

// Read reads data from a file source and indicates when the reading has completed
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

// GetModifiedTime returns the time when the source was last updated
func (r *Reader) GetModifiedTime() (time.Time, error) {

	info, err := os.Stat(r.source)

	if err != nil {
		return time.Unix(0, 0), err
	}

	return info.ModTime(), nil
}

// NewReader returns a new instance of Reader using the provided source identifier
func NewReader(source string) readers.Resource {

	return &Reader{source: source}
}
