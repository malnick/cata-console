package main

import (
	//"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	//"net/http"
	"os"
	"text/template"
)

type HostDashboard struct {
	Hostname string
}

// Accepts the Hostname and creates a new dashboard for the host in ./hostdata/templates/$hostname
func createHostDashboard(hostname string) {
	makeDirectories(hostname)
	// Init a new dashboard obj
	var hostdash HostDashboard
	// Update the hostname
	hostdash.Hostname = hostname
	// Parse a new json template and save it
	t, err := template.ParseFiles("templates/host_dashboard.json.template")
	if err != nil {
		log.Error(err)
	}
	err = t.Execute(os.Stdout, hostdash)
	if err != nil {
		log.Error(err)
	}
}

func updateHostDashboard() {

}

// Ensures we have a templates dir for this hostname
func makeDirectories(hostname string) {
	err := os.MkdirAll(fmt.Sprintf("hostdata/templates/%s", hostname), 0755)
	if err != nil {
		log.Error(err)
	}
}
