package main

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"gopkg.in/olivere/elastic.v2"
	"os"
	//"reflect"
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
		// Set index name
		index := strings.ToLower(host.Host)
		// Check index exists
		b, err := client.IndexExists(index).Do()
		if b == false {
			log.Debug("Creating ES Index ", index)
			_, err = client.CreateIndex(strings.ToLower(index)).Do()
			if err != nil {
				log.Warn("Index already created: ", index)
			}
		}

		host.Time = string(time.Now().Format("20060102150405"))
		//index := strings.Join([]string{strings.ToLower(host.Host), "-", host.Time}, "")

		// Dump the POST to elasticsearch, creating a new index based on timestamp data for the specific host
		_, err = client.Index().
			Index(index).
			Type(host.Host).
			Id("1").
			BodyJson(host).
			Do()
		if err != nil {
			log.Error("Problem dumping to ES: ", err)
		}
	}
}

func queryElastic(query string) (results HttpPost) {
	client, _ := elastic.NewClient()
	hostQuery := elastic.NewQueryStringQuery(query)
	searchResult, err := client.Search().
		Query(&hostQuery).
		Do()

	if err != nil {
		log.Error("Error during query: ", err)
	}

	if searchResult.Hits != nil {
		log.Info("Hits: ", searchResult.Hits.TotalHits)
		for _, hit := range searchResult.Hits.Hits {
			err := json.Unmarshal(*hit.Source, &results)
			if err != nil {
				log.Error(err)
			}
		}
	}
	return results
}
