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
	hostJsonFile := fmt.Sprintf("hostdata/templates/%s_dashbaord.json", hostname)
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
	f, err := os.Create(fmt.Sprintf(hostJsonFile, hostname))
	if err != nil {
		log.Error(err)
	}
	//Execute our template
	err = t.Execute(f, hostdash)
	if err != nil {
		log.Error(err)
	}
	updateHostDashboard(hostJsonFile)
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
	if err != nil {
		log.Error(err)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
	}
	defer resp.Body.Close()
	log.Info("POST to Grafana: ")
	fmt.Println(string(jsonFile))
	log.Info("Response: ", resp.Status)
}

// Ensures we have a templates dir for this hostname
func makeDirectories(hostname string) {
	err := os.MkdirAll(fmt.Sprintf("hostdata/templates/%s", hostname), 0755)
	if err != nil {
		log.Error(err)
	}
}
