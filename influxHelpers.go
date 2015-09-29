// Influx DB helpers

package main

import (
	"encoding/json"
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
	cmd := fmt.Sprintf("select * from /.*/ where hostname = '%s' limit 3", host)
	allData, err := queryInfluxDb(influxClient, cmd, InfluxDb)
	if err != nil {
		return allData, err
	}

	return allData, nil
}

// accepts []client.result and returns map[measurement][metric][]map[timestamp] = value
func transformResultsToMap(input []client.Result) (output map[string]map[string]interface{}) {
	// Accepts client.Result and maps it into a usable data structure for our pages
	log.Warn(input)

	output = make(map[string]map[string]interface{})

	// There is only a single index [0] returned from influx
	for _, v := range input {
		// v.Series is the tablature data
		for _, values := range v.Series {
			// Create a new output key from the value of the tablature data

			output[values.Name] = make(map[string]interface{})

			// For each index and metric column, range over
			for i, mc := range values.Columns {

				// For each metric index and metric values, range
				for _, mv := range values.Values {
					if mv[i] != nil {
						// Time is always the first in the values array
						timestamp := mv[0].(string)

						log.Warn("COLUMN ", mc)
						// Init a new map for the tablature name with a key for the metric
						output[values.Name][mc] = make([]map[string]interface{}, len(mv))

						//output[values.Name][mc].(map[string]interface{})[timestamp] = make([]map[string]string, len(mv))

						//outArry := make([]map[string]interface{}, len(mv))

						outMap := make(map[string]interface{})

						log.Warn(timestamp, ": ", mv[i])
						switch mv[i].(type) {
						case json.Number:
							outMap[timestamp] = mv[i].(json.Number)

							output[values.Name][mc].([]map[string]interface{}) = append(output[values.Name][mc].([]map[string]interface{}), outMap)

						case string:
							output[values.Name][mc].(map[string]interface{})[timestamp] = mv[i].(string)
						case uint8:
							output[values.Name][mc].(map[string]interface{})[timestamp] = mv[i].(uint8)
						}
					}
				}
			}
		}
	}
	for k, v := range output {
		log.Warn(k)
		for key, value := range v {
			log.Warn(key, " ", value)
		}
	}
	return output
}
