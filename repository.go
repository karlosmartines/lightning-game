package main

import (
	"log"

	uuid "github.com/satori/go.uuid"
)

type user struct {
	Email    string
	Password []byte
	Balance  int
}

/*type game struct {
	User user
	Bet string
	result string

}*/

var dbUsers = map[string]user{}
var dbSessions = map[string]string{}

func createUser(u user) string {
	uID := uuid.NewV4().String()
	dbUsers[uID] = u
	return uID
}

func readUser(id string) user {
	var u user
	u = dbUsers[id]
	return u
}
func updateUser(id string, u user) string {
	_, ok := dbUsers[id]
	if !ok {
		log.Fatalln("Update user not ok")
		return nil
	}
	dbUsers[id] = u
	return id
}
func deleteUser(id string) *user {
	u := findUser(id)
	if u == nil {
		return nil
	}
	delete(dbUsers, id)
	return u
}

func findUser(id string) *user {
	u, ok := dbUsers[id]
	if !ok {
		log.Fatalf("User %s not found", id)
		return nil
	}
	return &u
}
