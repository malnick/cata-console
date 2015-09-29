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
	LogLevel string  `json:"log_level"`
	Alarms   []Alarm `json:"alarms"`
}

func ParseEnv(c Config) Config {
	// Create a few matches for our env parsing down the road
	matchEnv, _ := regexp.Compile("CATA_ALARM_*")
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
	}
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
