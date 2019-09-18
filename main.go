package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"text/template"
	"time"

	uuid "github.com/satori/go.uuid"
)

type userAccount struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
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

var (
	tokens map[string]int
)

func main() {

	connectToDB()
	tokens = make(map[string]int)

	http.HandleFunc("/hello", helloHandler)
	http.HandleFunc("/library", auth(libraryHandler))
	http.HandleFunc("/addNote", addNoteHandler)
	http.HandleFunc("/changeStatus", changeStatusHandler)
	http.HandleFunc("/signin", signin)
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

func createUserAccountHandler(rw http.ResponseWriter, r *http.Request) {
	uA := &userAccount{}
	json.NewDecoder(r.Body).Decode(uA)
	err := createUserAccount(*uA)
	if err != nil {
		log.Println(err)
		rw.Write([]byte("something wrong"))
		return
	}
	rw.Write([]byte("success"))
}

func signin(rw http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("template/signin.html")
		if err != nil {
			log.Println(err)
			fmt.Fprint(rw, "something wrong")
		}
		tmpl.Execute(rw, nil)
	} else {
		supposedUA := userAccount{}
		err := json.NewDecoder(r.Body).Decode(&supposedUA)
		if err != nil {
			fmt.Println(err)
			rw.Write([]byte("badrequest"))
			return
		}
		currentUA, err := getUserAccount(supposedUA.Email)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		if checkPasswordHash(supposedUA.Password, currentUA.Password) {
			newSessionToken := uuid.Must(uuid.NewV4()).String()
			tokens[newSessionToken] = 0
			http.SetCookie(rw, &http.Cookie{
				Name:    "session_token",
				Value:   newSessionToken,
				Expires: time.Now().Add(time.Second * 120),
			})
		} else {
			rw.WriteHeader(http.StatusUnauthorized)
			rw.Write([]byte("wrong password"))
		}
	}
}

func auth(f http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("session_token")
		if err != nil {
			log.Println(err)
			rw.Write([]byte("who are you?"))
			return
		}
		if _, ok := tokens[c.Value]; ok {
			f(rw, r)
		} else {
			rw.Write([]byte("bad cookie"))
		}
	}
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
