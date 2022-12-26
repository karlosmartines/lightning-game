package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	uuid "github.com/satori/go.uuid"
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
