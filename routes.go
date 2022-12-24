package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		index,
	},
	Route{
		"Signup",
		"GET",
		"/signup",
		signup,
	},
	Route{
		"Login",
		"GET",
		"/login",
		login,
	},
	Route{
		"Logout",
		"GET",
		"logout",
		logout,
	},
	Route{
		"Play",
		"GET",
		"/play",
		play,
	},
}
