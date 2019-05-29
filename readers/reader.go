package readers

import (
	"time"
)

// Reader defines an interface that can read from and retrieve the last modified time of a data source
type Reader interface {
	Read(p []byte) (n int, err error)
	GetModifiedTime() (time.Time, error)
}
