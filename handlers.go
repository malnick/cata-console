package main

import (
	"encoding/json"
	//"fmt"
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"net/http"
)

type HttpPost struct {
	Host string
	Data map[string]interface{}
	Time string
}

// The container index
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
		dumpToElastic(hostDataArry)
	}

}

// The host index
func Console(w http.ResponseWriter, r *http.Request) {
	log.Debug("/ GET")
}
