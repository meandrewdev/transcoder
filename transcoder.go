package transcoder

import (
	"io"
)

// Transcoder ...
type Transcoder interface {
	Start(opts Options) ([]byte, error)
	Output(o string) Transcoder
	OutputPipe(w *io.WriteCloser, r *io.ReadCloser) Transcoder
	Probe(input string) (Metadata, error)
}
