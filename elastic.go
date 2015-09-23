package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/olivere/elastic.v2"
	"os"
	"strings"
	"time"
)

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

		// Dump the POST to elasticsearch, creating a new index based on timestamp data for the specific host
		_, err = client.Index().
			Index(index).
			BodyJson(host).
			Do()
		if err != nil {
			log.Error("Problem dumping to ES: ", err)
		}
	}
}

func queryElastic(query string) {
	client, _ := elastic.NewClient()
	hostQuery := elastic.NewTermQuery("Host", query)
	searchResult, err := client.Search().
		Query(&hostQuery). // specify the query
		Pretty(true).      // pretty print request and response JSON
		Do()               // execute
	if err != nil {
		// Handle error
		log.Error("Error querying ", hostQuery)
		panic(err)
	}
	log.Info(searchResult)
	if searchResult.Hits != nil {
		log.Info("Hits: ", searchResult.Hits.TotalHits)
		for _, hit := range searchResult.Hits.Hits {
			var h HttpPost
			err := json.Unmarshal(*hit.Source, &h)
			if err != nil {
				log.Error(err)
			}
			log.Info(h)
			log.Info(h.Data)
		}
	}
}
