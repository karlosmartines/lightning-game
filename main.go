package main

import (
	"log"
	"net/http"
	"text/template"

	"lightning-game/game"
)

var tpl *template.Template

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/play", play)
	http.HandleFunc("/login", logIn)
	http.ListenAndServe(":8080", nil)
}
func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}
func logIn(w http.ResponseWriter, r *http.Request) {
	welcomeMessage := r.FormValue("email")
	err := tpl.ExecuteTemplate(w, "index.html", welcomeMessage)
	if err != nil {
		log.Fatalln(err)
	}
}
func index(w http.ResponseWriter, r *http.Request) {
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
