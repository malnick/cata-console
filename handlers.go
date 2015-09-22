package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/olivere/elastic.v2"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type HttpPost struct {
	Host string
	Data map[string]interface{}
}

func dumpToElastic(data []byte) {
	// Create new elastic client
	client, err := elastic.NewClient()
	if err != nil {
		log.Debug("Failed to create elastic client")
		os.Exit(1)
	}

	// Create a new index for the host
	_, err = client.CreateIndex("twitter").Do()
	if err != nil {
	}

}

// The container index
func agent(w http.ResponseWriter, r *http.Request) {
	log.Debug("/agent POST")
	// Make a channel to dump our requests asynchronously
	respCh := make(chan *HttpPost)

	// Make an array of hostData to feed into
	hostDataArry := []*HttpPost{}

	// Spawn a proc to dump the data into our channel
	go func(r *http.Request) {
		var newData HttpPost
		body, err := ioutil.ReadAll(r.Body)
		// Unmarshal the POST into .Data
		err = json.Unmarshal(body, &newData.Data)
		// Type assert our way to the hostname
		newData.Host = newData.Data["host"].(map[string]interface{})["hostname"].(string)
		respCh <- &newData
	}(r)

	for count := 0; count != 10; count++ {
		select {
		case r := <-respCh:
			log.Debug("POST Received: ", r.Data)
			hostDataArry = append(hostDataArry, r)
		// Should count to 10
		case <-time.After(time.Second * 1):
			fmt.Printf(".")
		}
	}

	log.Debug("New Data:")
	fmt.Println(hostDataArry)
}

// The host index
func console(w http.ResponseWriter, r *http.Request) {
	log.Debug("/ GET")
}
