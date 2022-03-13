package main

import reader "go-stack-heap/2-ptr-iface-met-inline/reader"

func main() {
	b := make([]byte, 3)

	r := reader.New()

	n, err := r.Read(b)
	if err != nil {
		panic(err)
	}
	if n != 3 {
		panic("Expected read byte count to be 3")
	}
}
