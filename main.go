package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":8080", nil)
}

func helloHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "hello world")
}
