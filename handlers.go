package main

import (
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func index(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.html", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func game(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		err := tpl.ExecuteTemplate(w, "game.html", "")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
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
		un := r.FormValue("username")
		p, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.MinCost)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if usernameExists(un) {
			// http.Error(w, "Email is already used", http.StatusForbidden)
			tpl.ExecuteTemplate(w, "signup.html", "Username is already used")
			return
		}
		u := user{"", un, p, 0}
		createUser(&u)
		c := createSessionCookie(u)
		http.SetCookie(w, c)
		http.Redirect(w, r, "/game", http.StatusSeeOther)
		return
	}
}

func login(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/game", http.StatusSeeOther)
	}
	if r.Method == http.MethodPost {
		un := r.FormValue("username")
		p := r.FormValue("password")
		u, err := readUserByName(un)
		if err != nil {
			tpl.ExecuteTemplate(w, "index.html", true)
			return
		}
		err = bcrypt.CompareHashAndPassword(u.Password, []byte(p))
		if err != nil {
			tpl.ExecuteTemplate(w, "index.html", true)
			return
		}
		c := createSessionCookie(*u)
		http.SetCookie(w, c)
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
