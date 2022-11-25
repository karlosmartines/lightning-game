package main

import (
	"fmt"
	"log"
	"net/http"
	"text/template"

	"lightning-game/game"
)

var tpl *template.Template

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/playodd/", playOdd)
	http.HandleFunc("/playeven/", playEven)
	http.ListenAndServe(":8080", nil)
}
func init() {
	tpl = template.Must(template.ParseGlob("templates/*"))
}

func index(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		log.Fatalln(err)
	}
}
func playEven(w http.ResponseWriter, r *http.Request) {
	if game.Play(true) {
		fmt.Fprintf(w, "Win!")
	} else {
		fmt.Fprintf(w, "Loose!")
	}
}
func playOdd(w http.ResponseWriter, r *http.Request) {
	if game.Play(false) {
		fmt.Fprintf(w, "Win!")
	} else {
		fmt.Fprintf(w, "Loose!")
	}
}
