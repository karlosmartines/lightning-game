package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"

	"lightning-game/game"
)

type user struct {
	Email    string
	Password []byte
}

var tpl *template.Template
var dbUsers = map[string]user{}
var dbSessions = map[string]string{}

func main() {
	http.HandleFunc("/", index)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/play", play)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", logIn)
	http.ListenAndServe(":8080", nil)
}
func signup(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	if r.Method == http.MethodGet {
		tpl.ExecuteTemplate(w, "signup.html", "")
	}
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		p := r.FormValue("password")
		if _, ok := dbUsers[email]; ok {
			http.Error(w, "Email is already used", http.StatusForbidden)
			return
		}
		sID := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = email
		bs, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.MinCost)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		u := user{email, bs}
		dbUsers[email] = u
		http.Redirect(w, r, "/", http.StatusSeeOther)
		fmt.Printf("Email: %s Password: %s Session ID: %s", dbUsers[email].Email, dbUsers[email].Password, c.Value)
		return
	}
}
func alreadyLoggedIn(r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	email := dbSessions[c.Value]
	_, ok := dbUsers[email]
	return ok
}
func getUser(r *http.Request) user {
	var u user
	c, err := r.Cookie("session")
	if err != nil {
		return u
	}
	if email, ok := dbSessions[c.Value]; ok {
		u = dbUsers[email]
	}
	return u
}
func setCookie(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		sID := uuid.NewV4()
		c = &http.Cookie{
			Name:     "session",
			Value:    sID.String(),
			HttpOnly: true,
		}
		http.SetCookie(w, c)
	}
	var u user
	if un, ok := dbSessions[c.Value]; ok {
		u = dbUsers[un]
	}
	if r.Method == http.MethodPost {
		un := r.FormValue("email")
		p := r.FormValue("password")
		u = user{un, []byte(p)}
		dbSessions[c.Value] = un
		dbUsers[un] = u
		fmt.Printf("UN: %s P: %s", un, p)
	}
}
func readCookie(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("my-cookie")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
	}
	fmt.Println(c.String())
}
func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}
func logIn(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session")
	if err != nil {
		sID := uuid.NewV4()
		c = &http.Cookie{
			Name:     "session",
			Value:    sID.String(),
			HttpOnly: true,
		}
		http.SetCookie(w, c)
	}
	var u user
	if un, ok := dbSessions[c.Value]; ok {
		u = dbUsers[un]
	}
	if r.Method == http.MethodPost {
		un := r.FormValue("email")
		p := r.FormValue("password")
		u = user{un, []byte(p)}
		dbSessions[c.Value] = un
		dbUsers[un] = u
		fmt.Printf("UN: %s P: %s S: %s", un, p, c.Value)
	}
	tpl.ExecuteTemplate(w, "login.html", "")
}
func index(w http.ResponseWriter, r *http.Request) {
	setCookie(w, r)
	err := tpl.ExecuteTemplate(w, "index.html", "")
	if err != nil {
		log.Fatalln(err)
	}
}
func play(w http.ResponseWriter, r *http.Request) {
	var gameWon bool
	if r.FormValue("flexRadioDefault") == "playeven" {
		gameWon = game.Play(true)
	} else {
		gameWon = game.Play(false)
	}
	if gameWon {
		displayGameResult(w, "You won!")
	} else {
		displayGameResult(w, "You lost!")
	}
}
func displayGameResult(w http.ResponseWriter, result string) {
	err := tpl.ExecuteTemplate(w, "index.html", result)
	if err != nil {
		log.Fatalln(err)
	}
}
