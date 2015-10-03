package main

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"text/template"
)

type HostDashboard struct {
	Hostname string
}

// Accepts the Hostname and creates a new dashboard for the host in ./hostdata/templates/$hostname
func createHostDashboard(hostname string) {
	log.Debug("Creating new dashboard for ", hostname)
	c := ParseConfig()
	// get our homdir
	katahome := c.KataHome
	// Make our directories
	makeDirectories(hostname, katahome)
	// Get our host data file path
	hostJsonFile := fmt.Sprintf("%s/dashboard_templates/%s_dashboard.json", katahome, hostname)
	// If the host data file exists, don't do anything - else, create it and post to grafana
	//if _, err := os.Stat(hostJsonFile); os.IsNotExist(err) {
	log.Warn(hostJsonFile, " not found. Creating and executing new dashboard from template.")
	// Init a new dashboard obj
	var hostdash HostDashboard
	// Update the hostname
	hostdash.Hostname = hostname
	// Parse a new json template and save it
	t, err := template.ParseFiles("grafana_config/host_dashboard.json.template")
	if err != nil {
		log.Error("Issue parsing dashboard template for ", hostname)
		log.Error(err)
	}
	// Get a new file handle
	f, err := os.Create(hostJsonFile)
	if err != nil {
		log.Error("Issue creating host JSON file ", hostJsonFile)
		log.Error(err)
	}
	defer f.Close()
	//Execute our template
	err = t.Execute(f, hostdash)
	if err != nil {
		log.Error("Issue executing template ", hostdash)
		log.Error(err)
	}
	updateHostDashboard(hostJsonFile)
	//}
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
func makeDirectories(hostname string, katahome string) {
	log.Debug("Checking directories for ", fmt.Sprintf("%s/dashboard_templates", katahome))
	err := os.MkdirAll(fmt.Sprintf("%s/dashboard_templates", katahome), 0755)
	if err != nil {
		log.Error(err)
	}
}

func makeInfluxDatasource() {
	//c := ParseConfig()

}

func createGrafanaIframes(hostname string) (uris []string) {
	// Get a local config set to use the grafana uri and port
	c := ParseConfig()
	// Create a dashed hostname for the grafana host
	dashedHostname := strings.Replace(hostname, ".", "-", -1)
	grafanaUrl := c.GrafanaUrl
	// We have 10 IDed URIs for Grafana
	id := 1
	for id < 11 {
		newIframeUri := fmt.Sprintf("<iframe src=\"http://%s/dashboard-solo/db/%s?panelId=%s&fullscreen&from=now-15m&to=now\" id=\"graph%s\" width=\"500\" height=\"250\" frameborder=\"0\"></iframe>",
			grafanaUrl,
			dashedHostname,
			strconv.Itoa(id),
			strconv.Itoa(id))
		// Append the new URI to our returned array
		uris = append(uris, newIframeUri)
		log.Debug("New iframe ", newIframeUri)
		id += 1
	}
	return uris
}
