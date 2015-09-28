package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	//	"github.com/influxdb/influxdb/client"
	"html/template"
	"io/ioutil"
	"net/http"
)

type AllHostDataPage struct {
	Queries []HttpPost
	Host    string
}

type LatestHostDataPage struct {
	Data string
	Host string
}

type MainPage struct {
	AvailableHosts []string
	HostHits       int
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
	// Use a helper function to return all distinct hostnames from influx
	results, err := getUniqueHosts()
	// If things go wrong error but keep running
	if err != nil {
		log.Error(err)
	}
	// queryInflux retuns a []client.Result
	// for each query, range over and find the values of each series.
	// In the case of getUniqueHosts we have only a single return plus a timestamp
	// which we get by [0][1] and then type asserting our interface to a string and
	// appending that to our array to pass into the MainPage struct.
	for _, v := range results {
		log.Debug("Hosts: ", v.Series[0].Values[0][1])
		for _, host := range v.Series[0].Values[0][1].([]interface{}) {
			p.AvailableHosts = append(p.AvailableHosts, host.(string))
		}
	}
	t, _ := template.ParseFiles("views/MainPage.html")
	t.Execute(w, p)

}

func ConsoleHostnameLatest(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostname := vars["hostname"]
	results, err := getLatestHostData(hostname)
	if err != nil {
		log.Error(err)
	}

	var p LatestHostDataPage
	log.Warn("Latest data for ", hostname, ":")
	for _, v := range results {
		log.Warn(v.Series[0].Values[0])
		for _, values := range v.Series {
			log.Warn(values)
			for _, va := range values {
				log.Warn(va)
			}
		}
	}
	// Latest data struct
	// Execute template
	p.Host = hostname
	t, _ := template.ParseFiles("views/LatestHostData.html")
	t.Execute(w, p)
}

func ConsoleHostnameRoot(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hostname := vars["hostname"]
	results := queryHostnameAll(hostname)
	log.Debug("Queried all results for ", hostname)
	fmt.Println(results)
	// New data page
	var p AllHostDataPage
	p.Queries = results
	p.Host = vars["hostname"]
	log.Debug("RESULTS ", results)
	log.Debug(p)

	// Parse Template
	t, _ := template.ParseFiles("views/AllHostData.html")
	t.Execute(w, p)
}
