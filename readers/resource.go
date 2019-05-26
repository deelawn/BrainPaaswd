package readers

import (
	"time"
)

type Resource interface {
	Read(p []byte) (n int, err error)
	GetModifiedTime() (time.Time, error)
}