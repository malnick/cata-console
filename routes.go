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
		"ConsoleConfig",
		"GET",
		"/config",
		ConsoleConfig,
	},
	Route{
		"HostnameDashboardRoot",
		"GET",
		"/host/{hostname}/dashboard",
		ConsoleHostnameDashboardRoot,
	},
	Route{
		"HostnameDashboardMeasurement",
		"GET",
		"/host/{hostname}/dashboard/{measurement}",
		ConsoleHostnameDashboardMeasurement,
	},
	Route{
		"HostnameRoot",
		"GET",
		"/host/{hostname}",
		ConsoleHostnameRoot,
	},
	Route{
		"HostnameRootMeasurement",
		"GET",
		"/host/{hostname}/{measurement}",
		ConsoleHostnameRootMeasurement,
	},
	Route{
		"HostnameRootMeasurementMetric",
		"GET",
		"/host/{hostname}/{measurement}/{metric}",
		ConsoleHostnameRootMeasurementMetric,
	},
	Route{
		"HostnameRootMeasurementTimevalue",
		"GET",
		"/host/{hostname}/{measurement}/{metric}/{timestamp}",
		ConsoleHostnameRootMeasurementTimevalue,
	},
}
