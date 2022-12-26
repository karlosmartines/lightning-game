package main

import (
	"net/http"
	"text/template"

	"golang.org/x/crypto/bcrypt"
)

var tpl *template.Template

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

func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	bs, _ := bcrypt.GenerateFromPassword([]byte("asdf"), bcrypt.MinCost)
	createUser(user{"test@test.com", bs, 0})
}
