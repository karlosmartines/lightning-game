package main

import (
	"log"
	"net/http"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

func index(w http.ResponseWriter, r *http.Request) {
	setCookie(w, r)
	err := tpl.ExecuteTemplate(w, "index.html", "")
	if err != nil {
		log.Fatalln(err)
	}
}
func game(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		err := tpl.ExecuteTemplate(w, "game.html", "")
		if err != nil {
			log.Fatalln(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func play(w http.ResponseWriter, r *http.Request) {
	var victory bool
	if r.FormValue("flexRadioDefault") == "playeven" {
		victory = gameWon(true)
	} else {
		victory = gameWon(false)
	}
	if victory {
		displayGameResult(w, "You won!")
	} else {
		displayGameResult(w, "You lost!")
	}
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
			// http.Error(w, "Email is already used", http.StatusForbidden)
			tpl.ExecuteTemplate(w, "signup.html", "Email is already used")
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
		http.Redirect(w, r, "/game", http.StatusSeeOther)
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/game", http.StatusSeeOther)
	}
	if r.Method == http.MethodPost {
		email := r.FormValue("email")
		p := r.FormValue("password")
		u, ok := dbUsers[email]
		if !ok {
			// http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			tpl.ExecuteTemplate(w, "index.html", true)
			return
		}
		err := bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		if err != nil {
			// http.Error(w, "Username and/or password do not match", http.StatusForbidden)
			tpl.ExecuteTemplate(w, "index.html", true)
			return
		}
		sID := uuid.NewV4()
		c := &http.Cookie{
			Name:  "session",
			Value: sID.String(),
		}
		http.SetCookie(w, c)
		dbSessions[c.Value] = email
		http.Redirect(w, r, "/game", http.StatusSeeOther)
	}
	tpl.ExecuteTemplate(w, "index.html", "")
}

func logout(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		c, _ := r.Cookie("session")
		delete(dbSessions, c.Value)
		c = &http.Cookie{
			Name:   "session",
			Value:  "",
			MaxAge: -1,
		}
		http.SetCookie(w, c)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}
}
