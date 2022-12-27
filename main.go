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
func startMux() {
	http.HandleFunc("/", index)
	http.Handle("/favicon.ico", http.NotFoundHandler())
	http.HandleFunc("/game", game)
	http.HandleFunc("/play", play)
	http.HandleFunc("/signup", signup)
	http.HandleFunc("/login", login)
	http.HandleFunc("/logout", logout)
	http.ListenAndServe(":8080", nil)
}
func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
	bs, _ := bcrypt.GenerateFromPassword([]byte("asdf"), bcrypt.MinCost)
	createUser(&user{"", "user1", bs, 0})
}
