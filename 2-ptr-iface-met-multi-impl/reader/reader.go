package reader

import "io"

type readerV1 struct {
	b []byte
}

type readerV2 struct {
	b []byte
}

var _ io.Reader = &readerV1{}
var _ io.Reader = &readerV2{}

func New(vers string) io.Reader {
	if vers == "v1" {
		return &readerV1{
			b: []byte{1, 2, 3},
		}
	}

	return &readerV2{
		b: []byte{1, 2, 3},
	}
}

func (r *readerV1) Read(p []byte) (int, error) {
	n := copy(p, r.b)
	return n, nil
}

func (r *readerV2) Read(p []byte) (int, error) {
	n := copy(p, r.b)
	return n, nil
}
