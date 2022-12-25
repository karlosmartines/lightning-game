package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

func alreadyLoggedIn(r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	email := dbSessions[c.Value]
	_, ok := dbUsers[email]
	return ok
}
func displayGameResult(w http.ResponseWriter, result string) {
	err := tpl.ExecuteTemplate(w, "game.html", result)
	if err != nil {
		log.Fatalln(err)
	}
}

func gameWon(evenBet bool) bool {
	rand.Seed(time.Now().Unix())
	evenWin := rand.Intn(36)%2 == 0
	return evenWin == evenBet
}
