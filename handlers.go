package main

import (
	"encoding/json"
	//"fmt"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/olivere/elastic.v2"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
)

type HttpPost struct {
	Host string
	Data map[string]interface{}
	Time string
}

func dumpToElastic(data []*HttpPost) {
	// Create new elastic client
	client, err := elastic.NewClient()
	if err != nil {
		log.Debug("Failed to create elastic client")
		os.Exit(1)
	}

	// Create a new index for the host
	for _, host := range data {
		host.Time = string(time.Now().Format("20060102150405"))
		index := strings.Join([]string{strings.ToLower(host.Host), "-", host.Time}, "")
		//index := strings.ToLower(host.Host)
		//		b, err := client.IndexExists(index).Do()
		//		if b == false {
		//			log.Debug("Creating ES Index ", index, b)
		//			_, err = client.CreateIndex(strings.ToLower(index)).Do()
		//			if err != nil {
		//				log.Warn("Index already created: ", index)
		//			}
		//		}
		_, err = client.Index().
			Index(index).
			BodyJson(host).
			Do()
		if err != nil {
			log.Error("Problem dumping to ES: ", err)
		}
	}
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
