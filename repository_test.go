package main

import (
	"testing"
)

func TestCreateUser(t *testing.T) {
	tests := []user{
		{"testing@testing.net", "password", 3},
		{"test2@test2.net", "pass3", 100},
		{"TEST3@TEST3.NET", "PASS5", 100000},
	}
	for _, u := range tests {
		id := createUser()
		if id == "" || id == nil {
			t.Errorf("ID = %s", id)
		}
	}
}
