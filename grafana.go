package main

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"text/template"
)

type HostDashboard struct {
	Hostname string
}

// Accepts the Hostname and creates a new dashboard for the host in ./hostdata/templates/$hostname
func createHostDashboard(hostname string) {
	makeDirectories(hostname)
	// Get our host data file path
	hostJsonFile := fmt.Sprintf("hostdata/templates/%s_dashbaord.json", hostname)
	// If the host data file exists, don't do anything - else, create it and post to grafana
	if _, err := os.Stat(hostJsonFile); os.IsNotExist(err) {
		log.Info(hostJsonFile, " not found. Creating and executing new dashboard from template.")
		// Init a new dashboard obj
		var hostdash HostDashboard
		// Update the hostname
		hostdash.Hostname = hostname
		// Parse a new json template and save it
		t, err := template.ParseFiles("templates/host_dashboard.json.template")
		if err != nil {
			log.Error(err)
		}
		// Get a new file handle
		f, err := os.Create(hostJsonFile)
		if err != nil {
			log.Error(err)
		}
		defer f.Close()
		//Execute our template
		err = t.Execute(f, hostdash)
		if err != nil {
			log.Error(err)
		}
		updateHostDashboard(hostJsonFile)
	}
}

func updateHostDashboard(hostJsonFile string) {
	c := ParseConfig()
	url := fmt.Sprintf("http://%s/api/dashboards/db", c.GrafanaUrl)
	jsonFile, err := ioutil.ReadFile(hostJsonFile)
	if err != nil {
		log.Error(err)
	}
	log.Info("Updating Grafana Dashboard: ", url)
	// POST our data to grafana
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonFile))
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.GrafanaAuth))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	log.Info("Request Headers: ", req.Header)
	if err != nil {
		log.Error(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	log.Debug("POST to Grafana: ")
	log.Debug(string(jsonFile))
	log.Info("Response: ", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)
	log.Info("Response Body: ", string(body))
}

// Ensures we have a templates dir for this hostname
func makeDirectories(hostname string) {
	err := os.MkdirAll("hostdata/templates", 0755)
	if err != nil {
		log.Error(err)
	}
}
