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
	KataHome string  `json:"kata_home"`
	// Grafana configuration
	GrafanaUrl  string `json:"grafana_url"`
	GrafanaAuth string `json:"grafana_auth"`
	// Influx datasource config for grafana
	InfluxUrl               string `json:"influx_url"`
	InfluxPort              string `json:"influx_port"`
	InfluxUser              string `json:"influx_user"`
	InfluxPassword          string `json:"influx_password"`
	InfluxBasicAuthEnabled  string `json:"influx_basic_auth_enabled"`
	InfluxBasicAuthUser     string `json:"influx_basic_auth_user"`
	InfluxBasicAuthPassword string `json:"influx_basic_auth_password"`
}

// Default Configuration
var DefaultGrafanaUrl = "localhost:3000"
var DefaultGrafanaAuth = "KATA_GRAFANA_AUTH Not Set!"
var DefaultKataHome = fmt.Sprintf("%s/.kata", os.Getenv("HOME"))
var DefaultInfluxUrl = "localhost"
var DefaultInfluxPort = "8086"
var DefaultInfluxUser = "admin"
var DefaultInfluxPassword = "admin"
var DefaultInfluxBasicAuthEnabled = "false"
var DefaultInfluxBasicAuthUser = ""
var DefaultInfluxBasicAuthPassword = ""

func ParseEnv(c Config) Config {
	// Create a few matches for our env parsing down the road
	matchEnv, _ := regexp.Compile("KATA_ALARM_*")
	// Get the Grafana Console URL from ENV or set default
	matchGrafanaUrl, _ := regexp.Compile("KATA_GRAFANA_URL=*")
	// Get the auth bearer for grafana api
	matchGrafanaAuth, _ := regexp.Compile("KATA_GRAFANA_AUTH=*")
	// Get home kata
	matchKataHome, _ := regexp.Compile("KATA_HOME=*")
	// Get influx url
	matchInfluxUrl, _ := regexp.Compile("KATA_INFLUX_URL=*")
	// Get influx port
	matchInfluxPort, _ := regexp.Compile("KATA_INFLUX_PORT=*")
	// Get influx user
	matchInfluxUser, _ := regexp.Compile("KATA_INFLUX_USER=*")
	// Get influx password
	matchInfluxPassword, _ := regexp.Compile("KATA_INFLUX_PASSWORD=*")
	// Get influx basic auth enabled
	matchInfluxBasicAuthEnabled, _ := regexp.Compile("KATA_INFLUX_BASIC_AUTH_ENABLED=*")
	// Get influx basic auth user
	matchInfluxBasicAuthUser, _ := regexp.Compile("KATA_INFLUX_BASIC_AUTH_USER=*")
	// Get influx basic auth password
	matchInfluxBasicAuthPassword, _ := regexp.Compile("KATA_INFLUX_BASIC_AUTH_PASSWORD=*")

	// Set the defaults and override later
	c.GrafanaUrl = DefaultGrafanaUrl
	c.GrafanaAuth = DefaultGrafanaAuth
	c.KataHome = DefaultKataHome
	c.InfluxUrl = DefaultInfluxUrl
	c.InfluxPort = DefaultInfluxPort
	c.InfluxUser = DefaultInfluxUser
	c.InfluxPassword = DefaultInfluxPassword
	c.InfluxBasicAuthEnabled = DefaultInfluxBasicAuthEnabled
	c.InfluxBasicAuthUser = DefaultInfluxBasicAuthUser
	c.InfluxBasicAuthPassword = DefaultInfluxBasicAuthPassword

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
			c.GrafanaAuth = fmt.Sprintf("%s==", grafanaAuth)
		}
		if matchKataHome.MatchString(e) {
			newKataHome := strings.Split(e, "=")[1]
			c.KataHome = newKataHome
		}
		if matchInfluxUrl.MatchString(e) {
			newInfluxUrl := strings.Split(e, "=")[1]
			c.InfluxUrl = newInfluxUrl
		}
		if matchInfluxPort.MatchString(e) {
			newInfluxPort := strings.Split(e, "=")[1]
			c.InfluxPort = newInfluxPort
		}
		if matchInfluxUser.MatchString(e) {
			newInfluxUser := strings.Split(e, "=")[1]
			c.InfluxUser = newInfluxUser
		}
		if matchInfluxPassword.MatchString(e) {
			newInfluxPassword := strings.Split(e, "=")[1]
			c.InfluxPassword = newInfluxPassword
		}
		if matchInfluxBasicAuthEnabled.MatchString(e) {
			influxBasicAuthEnabled := strings.Split(e, "=")[1]
			c.InfluxBasicAuthEnabled = influxBasicAuthEnabled
		}
		if matchInfluxBasicAuthUser.MatchString(e) {
			influxBasicAuthUser := strings.Split(e, "=")[1]
			c.InfluxBasicAuthUser = influxBasicAuthUser
		}
		if matchInfluxBasicAuthPassword.MatchString(e) {
			influxBasicAuthPassword := strings.Split(e, "=")[1]
			c.InfluxBasicAuthPassword = influxBasicAuthPassword
		}
	}
	checkhome(c.KataHome)
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

func checkhome(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.MkdirAll(path, 0755)
		if err != nil {
			log.Error(err)
		}
	}
}
