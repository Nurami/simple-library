package main

import (
	"fmt"
	"net/http"
	"time"
)

type note struct {
	ID             string
	BookName       string
	EmergenceDate  time.Time
	BuyingDate     time.Time
	ReadingEndDate time.Time
	Status         string
}

func main() {
	http.HandleFunc("/hello", helloHandler)
	http.ListenAndServe(":8080", nil)
}

func helloHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "hello world")
}

func addNoteHandler(rw http.ResponseWriter, r *http.Request) {
}
func changeStatusHandler(rw http.ResponseWriter, r *http.Request) {
}
