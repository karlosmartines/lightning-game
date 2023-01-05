package main

import (
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type templateData struct {
	User *user
	Game *game
}

func index(w http.ResponseWriter, r *http.Request) {
	err := tpl.ExecuteTemplate(w, "index.html", "")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
func home(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		uID, err := getSessionUser(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		u, err := readUser(uID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		g := game{
			"",
			"",
			0,
			false,
			0,
		}
		td := templateData{
			u,
			&g,
		}
		err = tpl.ExecuteTemplate(w, "home.html", td)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
func play(w http.ResponseWriter, r *http.Request) {
	uID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := readUser(uID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	g := emptyGame(uID)
	td := templateData{
		u,
		g,
	}
	bs, _ := strconv.Atoi(r.FormValue("betsum"))
	if bs > u.Balance || bs < 1 || u.Balance < 1 {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
	g.Bet = bs
	bettype := r.FormValue("bettype")
	if bettype == "Even" {
		g.EvenBet = true
	}
	if gameWon(td.Game.EvenBet) {
		u.Balance += bs
		g.Result = bs * 2
		_, _ = updateUser(uID, *u)
	} else {
		u.Balance -= bs
		g.Result = -bs * 2
		_, _ = updateUser(uID, *u)
	}
	tpl.ExecuteTemplate(w, "home.html", td)
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
		http.Redirect(w, r, "/home", http.StatusSeeOther)
		return
	}
}

func fundAccount(w http.ResponseWriter, r *http.Request) {
	uID, err := getSessionUser(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u, err := readUser(uID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	u.Balance = 1000
	_, err = updateUser(uID, *u)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	nu, _ := readUser(uID)
	g := emptyGame("")
	td := templateData{
		nu,
		g,
	}
	tpl.ExecuteTemplate(w, "home.html", td)
}

func login(w http.ResponseWriter, r *http.Request) {
	if alreadyLoggedIn(r) {
		http.Redirect(w, r, "/home", http.StatusSeeOther)
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
		http.Redirect(w, r, "/home", http.StatusSeeOther)
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
