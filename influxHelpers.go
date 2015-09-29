// Influx DB helpers

package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/influxdb/influxdb/client"
)

// Accepts a an array of results and returns a singlular line query
func stringArrayResults(resultsArry []client.Result) []string {
	var stringArry = []string{}
	for _, v := range resultsArry {
		switch v.Series[0].Values[0][1].(type) {
		case []interface{}:
			for _, value := range v.Series[0].Values[0][1].([]interface{}) {
				switch value.(type) {
				case json.Number:
					stringArry = append(stringArry, string(value.(json.Number)))
				case string:
					stringArry = append(stringArry, value.(string))
				}
			}
		case json.Number:
			value := v.Series[0].Values[0][1].(json.Number)
			stringArry = append(stringArry, string(value))
		}
	}
	return stringArry
}

func getUniqueHosts() ([]string, error) {
	var uniqueHosts = []string{}
	log.Debug("Getting distinct hosts")
	// Get a fresh client
	influxClient := SetInflux()
	// Query influx for distinct hosts
	cmd := "select distinct(hostname) from host"
	distinctHosts, err := queryInfluxDb(influxClient, cmd, InfluxDb)
	if err != nil {
		return nil, err
	}
	// Get the output in []string format
	uniqueHosts = stringArrayResults(distinctHosts)
	// return the results
	return uniqueHosts, nil
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

// Counts all entries a host has in the DB
func countHostEntries(host string) string {
	log.Debug("Counting host entries for ", host)
	influxClient := SetInflux()
	cmd := fmt.Sprintf("select count(hostname) from /.*/ where hostname = '%s'", host)
	resp, err := queryInfluxDb(influxClient, cmd, InfluxDb)
	if err != nil {
		log.Error(err)
	}
	respArry := stringArrayResults(resp)
	return respArry[0]
}

// accepts []client.result and returns map[measurement][metric][]map[timestamp] = value
func transformResultsToMap(input []client.Result) (output map[string]map[string]interface{}) {
	// Accepts client.Result and maps it into a usable data structure for our pages
	log.Warn(input)
	// Init a new map with hte make func to create a new output for our downstream data
	output = make(map[string]map[string]interface{})
	// There is only a single index [0] returned from influx
	for _, v := range input {
		// v.Series is the tablature data
		for _, values := range v.Series {
			// Create a new output key from the value of the tablature data
			output[values.Name] = make(map[string]interface{})
			// For each index and metric column, range over
			for i, mc := range values.Columns {
				// outMap is our map of timestamped values for our interface
				outMap := make(map[string]interface{})
				// For each metric index and metric values, range
				for _, mv := range values.Values {
					if mv[i] != nil {
						// Time is always the first in the values array
						timestamp := mv[0].(string)
						// Init a new map for the tablature name with a key for the metric
						output[values.Name][mc] = make([]map[string]interface{}, len(mv))
						// Init a new array with length mv to dump our data to
						outArry := make([]map[string]interface{}, len(mv))
						// Debug output
						log.Debug(timestamp, ": ", mv[i])
						// Type assert our way into the mapped timestamp
						switch mv[i].(type) {
						case json.Number:
							outMap[timestamp] = mv[i].(json.Number)
							outArry = append(outArry, outMap)
						case string:
							outMap[timestamp] = mv[i].(string)
							outArry = append(outArry, outMap)
						case uint8:
							outMap[timestamp] = mv[i].(uint8)
							outArry = append(outArry, outMap)
						}
						// Append our timestamped values to the map
						output[values.Name][mc] = outMap

					}
				}
			}
		}
	}
	return output
}
