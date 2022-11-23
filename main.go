package main

import (
	"fmt"
	"net/http"

	"lightning-game/game"
)

func main() {
	http.HandleFunc("/", index)
	http.HandleFunc("/playodd/", playOdd)
	http.HandleFunc("/playeven/", playEven)
	http.ListenAndServe(":8080", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "hello, Karl")
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
