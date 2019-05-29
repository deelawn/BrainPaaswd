package file

import (
	"io"
	"os"
	"time"

	"github.com/deelawn/BrainPaaswd/readers"
)

// Reader is a file reader that implements io.Reader
type Reader struct {
	source string
	fd     *os.File
}

// Read reads data from a file source and indicates when the reading has completed
func (r *Reader) Read(p []byte) (n int, err error) {

	// Open the file if it is the first iteration of calling Read
	if r.fd == nil {
		r.fd, err = os.Open(r.source)

		if err != nil {
			return 0, err
		}
	}

	// Read from the file
	n, err = r.fd.Read(p)

	// Everything has been read; close the file
	if n == 0 && err == io.EOF {
		_ = r.fd.Close()
	}

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
func NewReader(source string) readers.Reader {

	return &Reader{source: source}
}
