package main

import reader "go-stack-heap/2-ptr-iface-met-multi-impl/reader"

func main() {
	b := make([]byte, 3)

	r := reader.New("v2")

	n, err := r.Read(b)
	if err != nil {
		panic(err)
	}
	if n != 3 {
		panic("Expected read byte count to be 3")
	}
}
