package main

import "github.com/influxdb/influxdb/client"
import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/url"
	//"os"
)

const (
	InfluxHost       = "localhost"
	InfluxPort       = 8086
	InfluxDb         = "hosts"
	HostMeasurements = "shapes"
)

func SetInflux() *client.Client {
	u, err := url.Parse(fmt.Sprintf("http://%s:%d", InfluxHost, InfluxPort))
	if err != nil {
		log.Error("Unable to parse url ", u)
		log.Fatal(err)
	}

	conf := client.Config{
		URL: *u,
		// Will have database auth later
		//Username: os.Getenv("INFLUX_USER"),
		//Password: os.Getenv("INFLUX_PWD"),
	}

	con, err := client.NewClient(conf)
	if err != nil {
		log.Error("Unable to create new client ", con)
		log.Fatal(err)
	}

	dur, ver, err := con.Ping()
	if err != nil {
		log.Error("Unable to run con.Ping() ", dur, " ", ver)
		log.Fatal(err)
	}
	log.Printf("Influx OK - response time %v, version %s", dur, ver)
	return con

}

func queryDB(con *client.Client, cmd string, db string) (res []client.Result, err error) {
	log.Debug("Influx Query: ", cmd)
	q := client.Query{
		Command:  cmd,
		Database: db,
	}
	if response, err := con.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	}
	return
}

func CheckDb(con *client.Client, db string) error {
	_, err := queryDB(con, fmt.Sprintf("CREATE DATABASE %s", db), db)
	if err != nil {
		return err
	}
	return nil
}

func influxify(data []*HttpPost) []client.Point {
	// Define our client point array
	var cpArry []client.Point
	// For each post in posts, dump the json to line protocol format
	for key, values := range data {
		log.Debug("Influxifying data for ", values.Host)
		log.Debug("Post ", key+1, " of ", len(data))
		// Parse the metrics and values
		// Build the dump into line protocol format
		// ex: memory,host=$hostname $memoryKey1=$memoryValue1,$memoryKey2=$memoryValue2 $timestamp
		for metricName, metricValues := range values.Data {
			// Define a new client point object and add data to it accordingly
			//var cp client.Point
			log.Debug(metricName)
			// Ensure our assertion only passes maps of strings and strings
			if _, ok := metricValues.(map[string]string); ok {
				for metrickey, measurement := range metricValues.(map[string]string) {
					log.Debug(fmt.Sprintf("%s: %s", metrickey, measurement))
				}
			}
		}
	}
	return cpArry
}

func dumpToInflux(host string, data []*HttpPost) (response string, err error) {
	log.Debug(fmt.Sprintf("%s: %s", host, data))
	// Create an influx client
	influxClient := SetInflux()
	// Check to make sure the DB is created, and create it if not
	err = CheckDb(influxClient, InfluxDb)
	if err != nil {
		log.Warn(fmt.Sprintf("Database %s: %s", InfluxDb, err))
	}
	// Influxify the JSON
	//influxData :=
	influxify(data)
	// Dump the data
	return "this", nil
}
