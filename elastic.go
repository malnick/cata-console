package main

import (
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
		//index := strings.ToLower(host.Host)
		//    b, err := client.IndexExists(index).Do()
		//    if b == false {
		//      log.Debug("Creating ES Index ", index, b)
		//      _, err = client.CreateIndex(strings.ToLower(index)).Do()
		//      if err != nil {
		//        log.Warn("Index already created: ", index)
		//      }
		//    }
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

}
