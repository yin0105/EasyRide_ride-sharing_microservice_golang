package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Distance struct {
	Text  string `json:"text"`
	Value string `json:"value"`
}

type Leg struct {
	Distance          Distance                 `json:"distance"`
	Duration          Distance                 `json:"duration"`
	EndAddress        string                   `json:"end_address"`
	EndLocation       map[string]string        `json:"end_location"`
	StartAddress      string                   `json:"start_address"`
	StartLocation     map[string]string        `json:"start_location"`
	Steps             []map[string]interface{} `json:"steps"`
	TrafficSpeedEntry []interface{}            `json:"traffic_speed_entry"`
	ViaWaypoint       []interface{}            `json:"via_waypoint"`
}

type Routes struct {
	Bounds           map[string]interface{} `json:"bounds"`
	Copyrights       string                 `json:"copyrights"`
	Legs             []Leg                  `json:"legs"`
	OverviewPolyline map[string]interface{} `json:"overview_polyline"`
	Summary          string                 `json:"summary"`
	Warnings         []interface{}          `json:"warnings"`
	WaypointOrder    []interface{}          `json:"waypoint_order"`
}

type HttpRes struct {

	// defining struct variables
	GeocodedWaypoints []interface{} `json:"geocoded_waypoints"`
	Routes            []Routes      `json:"routes"`
	Status            string        `json:"status"`
}

func main() {
	fmt.Println("Starting the application...")
	// response, err := http.Get("https://maps.googleapis.com/maps/api/directions/json?origin=37.75434337954133,%20-122.4837655029297&destination=137.750543040919084,%20122.41853417968751&key=AIzaSyDI57hkGB_K7Mtp4eFdYiy0mIw68z_1R1Y")
	response, err := http.Get("https://maps.googleapis.com/maps/api/directions/json?origin=37.75434337954133,%20-122.4837655029297&destination=37.750543040919084,%20-122.41853417968751&key=AIzaSyDI57hkGB_K7Mtp4eFdYiy0mIw68z_1R1Y")
	if err != nil {
		fmt.Printf("The HTTP request failed with error %s\n", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// fmt.Println(string(data))

		var result map[string]interface{}
		json.Unmarshal([]byte(data), &result)

		var result2 HttpRes
		json.Unmarshal(data, &result2)
		if len(result2.Routes) == 0 {
			fmt.Println("Routes == nil")
		} else {
			d := result2.Routes[0].Legs[0].Distance.Text
			d = strings.TrimSpace(d)
			unit := string(d[len(d)-2:])
			d = strings.TrimSpace(d[:len(d)-2])
			fmt.Println(unit)
			fmt.Println(d)
			if s, err := strconv.ParseFloat(d, 32); err == nil {
				if unit == "mi" {
					s *= 1.609
				}
				fmt.Println(s) // 3.1415927410125732
			}
		}

	}
}
