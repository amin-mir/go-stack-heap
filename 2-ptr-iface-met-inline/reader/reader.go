package reader

import "io"

type reader struct {
	b []byte
}

var _ io.Reader = &reader{}

func New() io.Reader {
	return &reader{
		b: []byte{1, 2, 3},
	}
}

func (r *reader) Read(p []byte) (int, error) {
	n := copy(p, r.b)
	return n, nil
}
