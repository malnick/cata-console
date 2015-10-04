// THe main Influx DB handlers.

package main

import "github.com/influxdb/influxdb/client"
import (
	"crypto/sha1"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/url"
	//"os"
	"time"
)

const (
	InfluxHost = "localhost"
	InfluxPort = 8086
	InfluxDb   = "hosts"
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
	log.Info(fmt.Sprintf("Influx OK - response time %v, version %s", dur, ver))
	return con

}

// Send commands to do things in influx, return the response.
func queryInfluxDb(con *client.Client, cmd string, db string) (res []client.Result, err error) {
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

// Query influx and see what we get back. Return an error if one exists.
func CheckDb(con *client.Client, db string) error {
	_, err := queryInfluxDb(con, fmt.Sprintf("CREATE DATABASE %s", db), db)
	if err != nil {
		return err
	}
	return nil
}

// Accepts a timestamp and returns the SHA1 string for use as a tag
func newShaStamp(t time.Time) string {
	stringThis := fmt.Sprintf("%s", t)
	hash := sha1.New()
	hash.Write([]byte(stringThis))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

// Accepts HTTP POST data and turns it into line protocol format
func influxify(data []*HttpPost) []client.Point {
	// Define our client point array
	var cpArry []client.Point
	// Get a sync timestamp
	timestamp := time.Now()
	shastamp := newShaStamp(timestamp)
	// For each post in posts, dump the json to line protocol format
	for key, values := range data {
		log.Debug("Influxifying data for ", values.Host)
		log.Debug("Post ", key+1, " of ", len(data))
		// Parse the metrics and values
		// Build the dump into line protocol format
		// ex: memory,host=$hostname $memoryKey1=$memoryValue1,$memoryKey2=$memoryValue2 $timestamp
		for metricName, metricValues := range values.Data {
			if metricValues == "null" {
				log.Warn("Metric values for ", metricName, " are null")
				break
			}
			// Define a new client point object and add data to it accordingly
			var cp client.Point
			// Initialize a new map of interfaces for our json from the agent
			cp.Fields = make(map[string]interface{})
			cp.Tags = make(map[string]string)
			// The metric we're about to influxify
			log.Debug(metricName)
			// Ensure our assertion type checks appropriately
			switch notSure := metricValues.(type) {
			case map[string]interface{}:
				for metrickey, measurement := range notSure {
					switch measurement.(type) {
					case float64:
						// Same for each block
						// Add timestamp
						cp.Time = timestamp
						// add the measurement
						cp.Measurement = metricName
						// add the measurement to the metrickey with correct type assertion
						log.Debug(fmt.Sprintf("%s: %s", metrickey, measurement))
						// When not in scientific notation, truncate
						cp.Fields[metrickey] = measurement.(float64)
						// tag it with the hostname for easy query later
						cp.Tags["hostname"] = values.Host
						cp.Tags["sha1"] = shastamp
					case int:
						cp.Time = timestamp
						cp.Measurement = metricName
						log.Debug(fmt.Sprintf("%s: %s", metrickey, measurement))
						cp.Fields[metrickey] = measurement.(int)
						cp.Tags["hostname"] = values.Host
						cp.Tags["sha1"] = shastamp
					case string:
						cp.Time = timestamp
						cp.Measurement = metricName
						log.Debug(fmt.Sprintf("%s: %s", metrickey, measurement))
						cp.Fields[metrickey] = measurement.(string)
						cp.Tags["hostname"] = values.Host
						cp.Tags["sha1"] = shastamp
					case int64:
						cp.Time = timestamp
						cp.Measurement = metricName
						log.Debug(fmt.Sprintf("%s: %s", metrickey, measurement))
						cp.Fields[metrickey] = measurement.(int64)
						cp.Tags["hostname"] = values.Host
						cp.Tags["sha1"] = shastamp
					case uint:
						cp.Time = timestamp
						cp.Measurement = metricName
						log.Debug(fmt.Sprintf("%s: %s", metrickey, measurement))
						cp.Fields[metrickey] = measurement.(uint)
						cp.Tags["hostname"] = values.Host
						cp.Tags["sha1"] = shastamp
					case uint64:
						cp.Time = timestamp
						cp.Measurement = metricName
						log.Debug(fmt.Sprintf("%s: %s", metrickey, measurement))
						cp.Fields[metrickey] = measurement.(uint64)
						cp.Tags["hostname"] = values.Host
						cp.Tags["sha1"] = shastamp
					}
				}
			// netcon and netio are both interface arrays
			case []interface{}:
				for _, block := range notSure {
					switch block.(type) {
					case map[string]interface{}:
						for metrickey, metricvalues := range block.(map[string]interface{}) {
							switch metricvalues.(type) {
							case float64:
								cp.Time = timestamp
								cp.Measurement = metricName
								log.Debug(fmt.Sprintf("%s: %s", metrickey, metricvalues))
								cp.Fields[metrickey] = metricvalues.(float64)
								cp.Tags["hostname"] = values.Host
								cp.Tags["sha1"] = shastamp
							case string:
								cp.Time = timestamp
								cp.Measurement = metricName
								log.Debug(fmt.Sprintf("%s: %s", metrickey, metricvalues))
								cp.Fields[metrickey] = metricvalues.(string)
								cp.Tags["hostname"] = values.Host
								cp.Tags["sha1"] = shastamp
							case int:
								cp.Time = timestamp
								cp.Measurement = metricName
								log.Debug(fmt.Sprintf("%s: %s", metrickey, metricvalues))
								cp.Fields[metrickey] = metricvalues.(int)
								cp.Tags["hostname"] = values.Host
								cp.Tags["sha1"] = shastamp
							case map[string]interface{}:
								cp.Time = timestamp
								cp.Measurement = metricName
								log.Debug(fmt.Sprintf("%s: %s", metrickey, metricvalues))
								cp.Fields[metrickey] = metricvalues.(map[string]interface{})
								cp.Tags["hostname"] = values.Host
								cp.Tags["sha1"] = shastamp
							}
						}
					}
				}
			default:
				log.Error("Could not match type for ", metricName)
			}
			// Let's check out data
			log.Info("New data for ", values.Host)
			log.Info("Tags: ", cp.Tags)
			log.Info("Measurement: ", cp.Measurement)
			log.Info("Fields: ", cp.Fields)
			log.Info("Timestamp: ", cp.Time)
			// Add our new point to the point arry
			cpArry = append(cpArry, cp)
		}
	}
	log.Debug("Full Dump: ", cpArry)
	return cpArry
}

// Accepts the hostname and the data from the JSON POST and sends that data to the
// influxifier to be dumped into influx in line protocol format
func dumpToInflux(host string, data []*HttpPost) (response string, err error) {
	log.Debug(fmt.Sprintf("%s: %s", host, data))
	// Create an influx client
	influxClient := SetInflux()
	// Check to make sure the DB is created, and create it if not
	err = CheckDb(influxClient, InfluxDb)
	if err != nil {
		log.Warn(fmt.Sprintf("Database %s: %s", InfluxDb, err))
	}
	// Dump the data
	batchDump := client.BatchPoints{
		Points:          influxify(data),
		Database:        InfluxDb,
		RetentionPolicy: "default",
		//Precision:       "s",
	}
	_, err = influxClient.Write(batchDump)
	if err != nil {
		log.Error("Could not dump data to ", InfluxDb)
		log.Error(err)
		return "Error", err
	}
	return "Success - data dumped to InfluxDB", nil
}
