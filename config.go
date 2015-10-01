package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Alarm struct {
	Name     string `json:"name"`
	Critical int    `json:"critical"`
	Warning  int    `json:"warning"`
	Ok       int    `json:"ok"`
}

type Config struct {
	LogLevel    string  `json:"log_level"`
	Alarms      []Alarm `json:"alarms"`
	GrafanaUrl  string  `json:"grafana_url"`
	GrafanaAuth string  `json:"grafana_auth"`
}

const (
	DefaultGrafanaUrl  = "localhost:3000"
	DefaultGrafanaAuth = "CATA_GRAFANA_AUTH Not Set!"
)

func ParseEnv(c Config) Config {
	// Create a few matches for our env parsing down the road
	matchEnv, _ := regexp.Compile("KATA_ALARM_*")
	// Get the Grafana Console URL from ENV or set default
	matchGrafanaUrl, _ := regexp.Compile("KATA_GRAFANA_URL=*")
	// Get the auth bearer for grafana api
	matchGrafanaAuth, _ := regexp.Compile("KATA_GRAFANA_AUTH=*")

	// Set the defaults and override later
	c.GrafanaUrl = DefaultGrafanaUrl
	c.GrafanaAuth = DefaultGrafanaAuth

	// Parse the env for our config
	for _, e := range os.Environ() {
		if matchEnv.MatchString(e) {
			var newAlarm Alarm
			log.Debug("New alarm found: ", e)
			newAlarm.Name = strings.Split(strings.Split(e, "=")[0], "CATA_ALARM_")[1]
			//	Crit value is first in list
			crit, _ := strconv.Atoi(strings.Split(strings.Split(e, "=")[1], ",")[0])
			warn, _ := strconv.Atoi(strings.Split(strings.Split(e, "=")[1], ",")[1])
			ok, _ := strconv.Atoi(strings.Split(strings.Split(e, "=")[1], ",")[2])
			newAlarm.Critical = crit
			newAlarm.Warning = warn
			newAlarm.Ok = ok
			c.Alarms = append(c.Alarms, newAlarm)
		}
		if matchGrafanaUrl.MatchString(e) {
			newGrafanaUrl := strings.Split(e, "=")[1]
			c.GrafanaUrl = newGrafanaUrl
		}
		if matchGrafanaAuth.MatchString(e) {
			grafanaAuth := strings.Split(e, "=")[1]
			c.GrafanaAuth = grafanaAuth
		}
	}
	// Plug the config into stdout so we have a record
	log.Info("Grafana URL: ", c.GrafanaUrl)
	log.Info("Grafana Auth: ", c.GrafanaAuth)
	// Get the consoles from the env
	return c
}

func ParseConfig() (c Config) {
	log.SetLevel(log.DebugLevel)
	if *verbose {
		log.Debug("Loglevel: Debug")
		c.LogLevel = "Verbose"
	} else if os.Getenv("VERBOSE") == "true" {
		log.Debug("LogLevel: Debug")
		c.LogLevel = "Verbose"
	} else {
		log.SetLevel(log.InfoLevel)
		log.Info("Loglevel: Info")
		c.LogLevel = "Info"
	}
	log.Info(fmt.Sprintf("Console running on :%s", *port))
	// Check influx connection
	_ = SetInflux()

	c = ParseEnv(c)
	return c
}
