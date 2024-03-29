package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
	//Alert       string
}

type Routes []Route

var routes = Routes{
	Route{
		"Agent",
		"POST",
		"/agent",
		Agent,
	},
	Route{
		"Console",
		"GET",
		"/",
		Console,
	},
	Route{
		"HostnameLatest",
		"GET",
		"/host/{hostname}/latest",
		ConsoleHostnameLatest,
	},
	Route{
		"HostnameRoot",
		"GET",
		"/host/{hostname}",
		ConsoleHostnameRoot,
	},
}
