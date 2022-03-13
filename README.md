# go-stack-heap

I'm assuming you know about Stack and Heap while reading this.
A good resource on this topic is this [video](https://www.youtube.com/watch?v=ZMZpH4yT7M0) from GopherCon SG 2019.

We'll focus more on escape analysis in Golang and ways to find whether a variable
would be allocated on the heap or stack. In short, we're going to run the following
command in each directory and reason about the output:

```bash
go build -gcflags="-m" .

# more verbose
go build -gcflags="-m -m" .
```

## Some random & seemingly irrelevant questions

We want to answer these questions by the end of README!

### HTTP handlers
```golang
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	...
})
```

Why use `http.ResponseWriter` as the first argument to the function even though techincally
it's for returning the output of the API? One could as well change it like so:
```golang
http.HandleFunc("/", func(r *http.Request) *http.Response {
	...
})
```

### gRPC handlers
The handler signature follows the pattern below ([link to example](https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_server/main.go)):

```golang
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	...
}
```

But there are other open source libraries such as `go-micro` which tweak it a little bit as follows ([link to example](https://github.com/asim/go-micro/blob/master/examples/helloworld/main.go)):

```golang
func (s *server) SayHello(ctx context.Context, req *pb.HelloRequest, rsp *pb.HelloReply) error {
	...
}
```

### Personal opinion/style
I think in most cases, we should accept values/pointers and return values unless the struct 
contains a huge amount of data. In such cases, we can just return pointers to structs when we don't
really care about the GC/heap allocations or pass a pointer to struct to our function argument and
fill it in the body.

```golang
func process(arg *arg) res {
    ...
}

func processperf(arg *arg, res *res) {
    ...
}
```

## Returning references to local memory of a function causes heap allocations

[code](1-return-pointer-heap)

```
# go-stack-heap/1-return-pointer-stack
./main.go:14:6: can inline returnResult
./main.go:3:6: can inline main
./main.go:4:21: inlining call to returnResult
./main.go:4:21: &result{...} does not escape
./main.go:6:8: "Unexpected output" escapes to heap
./main.go:15:9: &result{...} escapes to heap        <== *****
```

However, if we return values ([code](1-return-value-stack)):

```
# go-stack-heap/1-return-value-stack
./main.go:3:6: can inline main
./main.go:6:8: "Unexpected output" escapes to heap
```

## Passing a pointer to a method of interface

[code](2-ptr-iface-met-inline)

```
# go-stack-heap/2-ptr-iface-met-inline
./main.go:8:17: inlining call to reader.New
./main.go:10:18: devirtualizing r.Read to *reader.reader
./main.go:6:11: make([]byte, 3) does not escape                        <== *****
./main.go:8:17: &reader.reader{...} does not escape
./main.go:8:17: []byte{...} does not escape                            <== *****
```

The buffer does not escape to the heap, however it's only because the compiler is able
to deduce that `io.Reader` is the same as `*reader.Reader` (devirtualizing). If we make
it a bit more complex so that there are several implementations of the same interface
then the compiler cannot make this optimization ([code](2-ptr-iface-met-multi-impl)).

```
# go-stack-heap/2-ptr-iface-met-multi-impl
./main.go:8:17: inlining call to reader.New
./main.go:6:11: make([]byte, 3) escapes to heap                            <== *****
./main.go:8:17: &reader.readerV1{...} escapes to heap
./main.go:8:17: []byte{...} escapes to heap                                <== *****
./main.go:8:17: &reader.readerV2{...} escapes to heap
./main.go:8:17: []byte{...} escapes to heap                                <== *****
```

## Escape Analysis for go-grpc handlers
Running `go build -gcflags="-m" .` for the example [here](https://github.com/grpc/grpc-go/tree/master/examples/helloworld/greeter_server) results in the following output:

```
./main.go:44:39: inlining call to helloworld.(*HelloRequest).GetName
./main.go:45:54: inlining call to helloworld.(*HelloRequest).GetName
./main.go:49:12: inlining call to flag.Parse
./main.go:55:26: inlining call to helloworld.RegisterGreeterServer
./main.go:34:2: can inline init
./main.go:34:17: inlining call to flag.Int
./main.go:55:26: devirtualizing helloworld.s.RegisterService to *grpc.Server
./main.go:43:7: s does not escape
./main.go:43:27: ctx does not escape
./main.go:43:48: leaking param content: in
./main.go:44:12: ... argument does not escape
./main.go:44:39: string(~R0) escapes to heap
./main.go:45:9: &helloworld.HelloReply{...} escapes to heap                     <== *****
./main.go:45:42: "Hello " + string(~R0) escapes to heap
./main.go:50:43: ... argument does not escape
./main.go:50:51: *port escapes to heap
./main.go:52:13: ... argument does not escape
./main.go:55:30: &server{} escapes to heap
./main.go:56:12: ... argument does not escape
./main.go:58:13: ... argument does not escape
<autogenerated>:1: .this does not escape
```

## Escape Analysis for go-micro handlers
Running `go build -gcflags="-m" .` for the example [here](https://github.com/asim/go-micro/tree/master/examples/helloworld) results in the following output:

```
./main.go:13:6: can inline (*Greeter).Hello
./main.go:19:29: inlining call to micro.NewService
./main.go:13:7: g does not escape
./main.go:13:25: ctx does not escape
./main.go:13:46: req does not escape
./main.go:13:63: rsp does not escape                  <== *****
./main.go:14:26: "Hello " + req.Name escapes to heap  <== *****
./main.go:19:29: []micro.Option{...} does not escape
./main.go:25:49: new(Greeter) escapes to heap
./main.go:28:12: ... argument does not escape
```
