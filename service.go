package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"
)

func getSessionUser(r *http.Request) (string, error) {
	c, err := r.Cookie("session")
	if err != nil {
		return "", err
	}
	s, err := readSession(c.Value)
	if err != nil {
		return "", err
	}
	return s.User, nil
}

func alreadyLoggedIn(r *http.Request) bool {
	c, err := r.Cookie("session")
	if err != nil {
		return false
	}
	s := dbSessions[c.Value]
	_, ok := dbUsers[s.User]
	return ok
}
func usernameExists(e string) bool {
	for _, u := range dbUsers {
		if u.Username == e {
			return true
		}
	}
	return false
}
func createSessionCookie(u user) *http.Cookie {
	s := createSession(u)
	c := &http.Cookie{
		Name:  "session",
		Value: s.Id,
	}
	return c
}

/*
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
			u = user{un, []byte(p), 0}
			dbSessions[c.Value] = un
			dbUsers[un] = u
			fmt.Printf("UN: %s P: %s", un, p)
		}
	}
*/
func displayGameResult(w http.ResponseWriter, result string) {
	err := tpl.ExecuteTemplate(w, "home.html", result)
	if err != nil {
		log.Fatalln(err)
	}
}

func gameWon(evenBet bool) bool {
	rand.Seed(time.Now().Unix())
	evenWin := rand.Intn(36)%2 == 0
	return evenWin == evenBet
}
