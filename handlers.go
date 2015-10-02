package main

import (
	"encoding/json"
	//	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	//	"github.com/influxdb/influxdb/client"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type AllHostDataPage struct {
	Queries map[string]map[string]interface{}
	Host    string
}

type LatestHostDataPage struct {
	Data       map[string]map[string]interface{}
	Host       string
	DashedHost string
}

type MainPage struct {
	AvailableHosts map[string]string
}

type HttpPost struct {
	Host string
	Data map[string]interface{}
	Time string
}

// /agent route - handles POSTs from aggregated agents
// dumps the POST to influxDb
func Agent(w http.ResponseWriter, r *http.Request) {
	log.Debug("/agent POST")
	// Make a channel to dump our requests asynchronously
	respCh := make(chan *HttpPost)

	// Make an array of hostData to feed into
	hostDataArry := []*HttpPost{}

	// Spawn a proc to dump the data into our channel
	go func(r *http.Request) {
		var newData HttpPost
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.Error(err)
		}
		// Unmarshal the POST into .Data
		err = json.Unmarshal(body, &newData.Data)
		// Type assert our way to the hostname
		newData.Host = newData.Data["host"].(map[string]interface{})["hostname"].(string)
		//newData.Time = string(time.Now().Format("2006010215040500"))
		respCh <- &newData
	}(r)

	// Check the channel for a resp
	select {
	case r := <-respCh:
		//	log.Debug("New data from ", r.Host, "@", r.Time)
		log.Debug("New data from ", r.Host)
		log.Debug(r.Data)
		hostDataArry = append(hostDataArry, r)
		//		dumpToElastic(hostDataArry)
		dumpToInflux(r.Host, hostDataArry)
	}
}

// The console root index /
// main console route:
//   - Displays known hosts, and basic info about them such as alarms
//   or other unique, top level data.
func Console(w http.ResponseWriter, r *http.Request) {
	log.Debug("/ GET")
	// Get a local main page struct to dump our data to
	var p MainPage
	// Init a new map
	p.AvailableHosts = make(map[string]string)
	// Use a helper function to return all distinct hostnames from influx
	uniqueHosts, err := getUniqueHosts()
	// If things go wrong error but keep running
	if err != nil {
		log.Error(err)
	}
	// For each unique host, count the number of host entries
	for _, host := range uniqueHosts {
		p.AvailableHosts[host] = countHostEntries(host)
	}
	t, _ := template.ParseFiles("views/MainPage.html")
	t.Execute(w, p)

}

func ConsoleHostnameLatest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostname := vars["hostname"]
	// Create grafana dashboard for our hostname
	log.Info("Request for latest host data for ", hostname)
	// Creates new json template for hostname and POSTs it to Grafana if it doesn't exist
	createHostDashboard(hostname)
	// Dash the hostname - Grafana URLs turn '.' into '-' so hostnames need to be sanitized
	dashHostname := strings.Replace(hostname, ".", "-", -1)

	// Get a local Latest Host Data strcut and init a new map for the
	var p LatestHostDataPage
	results, err := getLatestHostData(hostname)
	if err != nil {
		log.Error(err)
	}
	// Add our dashed hostname parameter
	p.DashedHost = dashHostname

	// Take the results from influx and transform them into something other than influx line protocol format
	p.Data = transformResultsToMap(results)

	log.Debug("Latest data for ", hostname, ":")

	// Execute template
	p.Host = hostname
	t, _ := template.ParseFiles("views/LatestHostData.html")
	t.Execute(w, p)
}

func ConsoleHostnameRoot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostname := vars["hostname"]
	// Query all host data
	results, _ := getAllHostData(hostname)
	log.Debug("Queried all results for ", hostname)
	// New data page
	var p AllHostDataPage
	mapped := transformResultsToMap(results)

	p.Queries = mapped

	p.Host = vars["hostname"]

	// Parse Template
	t, _ := template.ParseFiles("views/AllHostData.html")
	t.Execute(w, p)
}

func ConsoleHostnameRootMeasurement(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Get hostname from passed URL
	hostname := vars["hostname"]
	// Get measurement from passed URL
	measurement := vars["measurement"]

	results, err := getAllHostDataMeasure(hostname, measurement)
	if err != nil {
		log.Error(err)
	}
	var p AllHostDataPage
	mapped := transformResultsToMap(results)

	p.Queries = mapped

	p.Host = vars["hostname"]

	// Parse Template
	t, _ := template.ParseFiles("views/AllHostData.html")
	t.Execute(w, p)
}

func ConsoleHostnameRootMeasurementTimevalue(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	// Get hostname from passed URL
	hostname := vars["hostname"]
	// Get measurement from passed URL
	measurement := vars["measurement"]
	// Get sanitized time
	timestamp := vars["timestamp"]

	results, err := getTimevalueHostDataMeasure(hostname, measurement, timestamp)
	if err != nil {
		log.Error(err)
	}
	var p AllHostDataPage
	mapped := transformResultsToMap(results)

	p.Queries = mapped

	p.Host = vars["hostname"]

	// Parse Template
	t, _ := template.ParseFiles("views/MeasurementByTimeHostData.html")
	t.Execute(w, p)
}

// Config page
func ConsoleConfig(w http.ResponseWriter, r *http.Request) {
	c := ParseConfig()
	t, _ := template.ParseFiles("views/Configuration.html")
	t.Execute(w, c)
}
