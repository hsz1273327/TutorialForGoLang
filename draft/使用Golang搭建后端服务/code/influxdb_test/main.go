package main

import (
	"fmt"
	"time"

	_ "github.com/influxdata/influxdb1-client" // this is important because of the bug in go mod
	influx "github.com/influxdata/influxdb1-client/v2"
)

func randomWrite() {
	client, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
		return
	}
	defer client.Close()
	bp, _ := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database:  "BumbleBeeTuna",
		Precision: "s",
	})
	tags := map[string]string{"cpu": "cpu-total"}
	fields := map[string]interface{}{
		"idle":   10.1,
		"system": 53.3,
		"user":   46.6,
	}
	pt, err := influx.NewPoint("cpu_usage", tags, fields, time.Now())
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
	bp.AddPoint(pt)
	err = client.Write(bp)
	if err != nil {
		fmt.Println("Error: ", err.Error())
	}
}
func createdb() {
	client, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
		return
	}
	defer client.Close()
	q := influx.NewQuery("CREATE DATABASE BumbleBeeTuna", "", "")
	if response, err := client.Query(q); err == nil && response.Error() == nil {
		fmt.Println("create db get:", response.Results)
	}
}

func query_point() {
	client, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr: "http://localhost:8086",
	})
	if err != nil {
		fmt.Println("Error creating InfluxDB Client: ", err.Error())
		return
	}
	defer client.Close()
	q := influx.NewQuery("SELECT count(*) FROM cpu_usage", "BumbleBeeTuna", "s")
	if response, err := client.Query(q); err == nil && response.Error() == nil {
		fmt.Println(response.Results)
	}
}

func main() {
	createdb()
	randomWrite()
	query_point()
}
