package main

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type user struct {
	Id       string
	Username string
	Password []byte
	Balance  int
}

type game struct {
	Id      string
	User    string
	Bet     int
	EvenBet bool
	Result  int
}

type appSession struct {
	Id   string
	User string
}

var dbUsers = map[string]user{}
var dbSessions = map[string]appSession{}
var dbGames = map[string]game{}

func createSession(u user) *appSession {
	s := appSession{
		uuid.NewV4().String(),
		u.Id,
	}
	dbSessions[s.Id] = s
	return &s
}

func readSession(sID string) (*appSession, error) {
	s, ok := dbSessions[sID]
	if !ok {
		return nil, fmt.Errorf("did not find session by id %s", s.Id)
	}
	return &s, nil
}

func createUser(u *user) *user {
	u.Id = uuid.NewV4().String()
	dbUsers[u.Id] = *u
	return u
}
func readUserByName(n string) (*user, error) {
	for _, u := range dbUsers {
		if u.Username == n {
			return &u, nil
		}
	}
	return nil, fmt.Errorf("username %s not found", n)
}
func readUser(id string) (*user, error) {
	u, ok := dbUsers[id]
	if !ok {
		return nil, fmt.Errorf("did not find user by id %s", id)
	}
	return &u, nil
}
func updateUser(id string, u user) (string, error) {
	_, ok := dbUsers[id]
	if !ok {
		return "", fmt.Errorf("update user %s returned ok = %t", id, ok)
	}
	dbUsers[id] = u
	return id, nil
}
func deleteUser(id string) *user {
	u, err := readUser(id)
	if err != nil {
		return nil
	}
	delete(dbUsers, id)
	return u
}
