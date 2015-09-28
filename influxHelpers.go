// Influx DB helpers

package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/influxdb/influxdb/client"
)

func getUniqueHosts() ([]client.Result, error) {
	log.Debug("Getting distinct hosts")
	// Get a fresh client
	influxClient := SetInflux()
	// Query influx for distinct hosts
	cmd := "select distinct(hostname) from host"
	distinctHosts, err := queryInfluxDb(influxClient, cmd, InfluxDb)
	if err != nil {
		return nil, err
	}
	// return the results
	return distinctHosts, nil
}

func getLatestHostData(host string) ([]client.Result, error) {
	log.Debug("Getting latest data for host ", host)
	// Get a new client
	influxClient := SetInflux()
	// Create the cmd to get latest data for host
	cmd := fmt.Sprintf("select * from /.*/ where hostname = '%s' limit 1", host)
	latestData, err := queryInfluxDb(influxClient, cmd, InfluxDb)
	if err != nil {
		return latestData, err
	}
	return latestData, nil
}

func getAllHostData(host string) ([]client.Result, error) {
	log.Debug("Getting all host data for ", host)
	// Get the new client
	influxClient := SetInflux()
	// Cmd to query all data for host
	cmd := fmt.Sprintf("select * from /.*/ where hostname = '%s'", host)
	allData, err := queryInfluxDb(influxClient, cmd, InfluxDb)
	if err != nil {
		return allData, err
	}

	return allData, nil
}

func transformResultsToMap(input []client.Result) (output map[string]map[string]interface{}) {
	// Accepts client.Result and maps it into a usable data structure for our pages
	output = make(map[string]map[string]interface{})
	for _, v := range input {
		for _, values := range v.Series {
			output[values.Name] = make(map[string]interface{})
			for i, mc := range values.Columns {
				if values.Values[0][i] != nil {
					output[values.Name][mc] = values.Values[0][i]
				}
			}
		}
	}
	return output
}
