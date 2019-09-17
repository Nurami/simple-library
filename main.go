package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"
)

type userAccount struct {
	ID       int
	Name     string
	Email    string
	Password string
}
type note struct {
	ID             int
	BookName       string
	EmergenceDate  time.Time
	BuyingDate     time.Time
	ReadingEndDate time.Time
	Status         string
}
type notesByStatus struct {
	WantNotes   []note
	BoughtNotes []note
	ReadNotes   []note
}

func main() {

	connectToDB()

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/library", libraryHandler)
	http.HandleFunc("/addNote", addNoteHandler)
	http.HandleFunc("/changeStatus", changeStatusHandler)
	http.ListenAndServe(":8080", nil)
}

func helloHandler(rw http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(rw, "hello world")
}

func libraryHandler(rw http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("template/index.html")
	if err != nil {
		log.Println(err)
		fmt.Fprint(rw, "something wrong")
	}
	notes := getNotes()
	nbs := notesByStatus{}
	nbs.create(notes)
	tmpl.Execute(rw, nbs)
}

func addNoteHandler(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}
	bookName := string(body)
	message := "Name of book is empty"
	if bookName != "" {
		message = "success"
		note := note{BookName: bookName, EmergenceDate: time.Now(), Status: "want"}
		err = addNote(note)
		if err != nil {
			log.Println(err)
		}
	}
	rw.Write([]byte(message))
}

func changeStatusHandler(rw http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Panic(err)
	}
	ID := string(body)
	message := "Something wrong"
	if ID != "" {
		message = "success"
		note := getNoteByID(ID)
		err = changeDateAndStatusOfNote(note)
		if err != nil {
			log.Println(err)
		}
	}
	rw.Write([]byte(message))
}

func (nbs *notesByStatus) create(notes []note) {
	for _, val := range notes {
		switch val.Status {
		case "want":
			nbs.WantNotes = append(nbs.WantNotes, val)
		case "bought":
			nbs.BoughtNotes = append(nbs.BoughtNotes, val)
		case "read":
			nbs.ReadNotes = append(nbs.ReadNotes, val)
		}
	}
}
