package main

import (
	"fmt"
	"net/http"
	"text/template"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Email    string
	Password []byte
}

var tpl *template.Template
var dbUsers = map[string]user{}
var dbSessions = map[string]string{}

func main() {
	startMux()
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

/*
	func readCookie(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("my-cookie")
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
		}
		fmt.Println(c.String())
	}
*/
func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	bs, _ := bcrypt.GenerateFromPassword([]byte("asdf"), bcrypt.MinCost)
	dbUsers["test@test.com"] = user{"test@test.com", bs}
}
