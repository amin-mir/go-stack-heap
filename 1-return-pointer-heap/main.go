package main

func main() {
	res := returnResult()
	if res.count != 10 {
		panic("Unexpected output")
	}
}

type result struct {
	count int
}

func returnResult() *result {
	return &result{
		count: 10,
	}
}
