package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, *content)
}

var (
	port    *string
	content *string
)

// CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build simple.go
func main() {

	port = flag.String("port", "8080", "a string")
	content = flag.String("word", "hello", "a string")
	flag.Parse()
	http.HandleFunc("/", helloHandler)
	portstr := fmt.Sprintf(":%s", *port)
	if err := http.ListenAndServe(portstr, nil); err != nil {
		panic(err)
	}
}
