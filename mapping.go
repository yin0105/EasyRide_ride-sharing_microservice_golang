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

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Step struct {
	Distance         Distance    `json:"distance"`
	Duration         Distance    `json:"distance"`
	EndLocation      Location    `json:"distance"`
	HtmlInstructions string      `json:"distance"`
	Polyline         interface{} `json:"distance"`
	StartLocation    Location    `json:"distance"`
	TravelMode       string      `json:"distance"`
}

type Leg struct {
	Distance          Distance      `json:"distance"`
	Duration          Distance      `json:"duration"`
	EndAddress        string        `json:"end_address"`
	EndLocation       Location      `json:"end_location"`
	StartAddress      string        `json:"start_address"`
	StartLocation     Location      `json:"start_location"`
	Steps             []Step        `json:"steps"`
	TrafficSpeedEntry []interface{} `json:"traffic_speed_entry"`
	ViaWaypoint       []interface{} `json:"via_waypoint"`
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

type MapRes struct {
	Distance string `json:"distance"`
	ARoad    int    `json:"a_road"`
}

func convertKm(d string) (float64, error) {
	d = strings.TrimSpace(d)
	unit := string(d[len(d)-2:])
	d = strings.TrimSpace(d[:len(d)-2])
	if s, err := strconv.ParseFloat(d, 32); err == nil {
		if unit == "mi" {
			s *= 1.609
		} else if unit == "ft" {
			s *= 0.0003048
		}
		return s, nil
	} else {
		return 0, err
	}
}

func handleMapping(w http.ResponseWriter, r *http.Request) {
	var mapRes MapRes
	mapRes.ARoad = 0
	mapRes.Distance = "0"
	fmt.Println("Starting Mapping Microservice...")
	uri := strings.Split(r.URL.Path, "/api/v1/mapping/")
	fmt.Println(uri[1])
	// response, err := http.Get("https://maps.googleapis.com/maps/api/directions/json?origin=37.75434337954133,%20-122.4837655029297&destination=137.750543040919084,%20122.41853417968751&key=AIzaSyDI57hkGB_K7Mtp4eFdYiy0mIw68z_1R1Y")
	response, err := http.Get("https://maps.googleapis.com/maps/api/directions/json?key=AIzaSyDI57hkGB_K7Mtp4eFdYiy0mIw68z_1R1Y&" + uri[1]) //origin=37.75434337954133,%20-122.4837655029297&destination=37.750543040919084,%20-122.41853417968751")
	fmt.Println("https://maps.googleapis.com/maps/api/directions/json?key=AIzaSyDI57hkGB_K7Mtp4eFdYiy0mIw68z_1R1Y&" + uri[1])
	if err != nil {
		http.Error(w, fmt.Sprintf("The HTTP request failed with error %s\n", err), http.StatusBadRequest)
		return
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		fmt.Println("data = " + string(data))
		var result map[string]interface{}
		json.Unmarshal([]byte(data), &result)

		fmt.Println(result)

		var result2 HttpRes
		json.Unmarshal(data, &result2)
		if len(result2.Routes) == 0 {
			fmt.Println("Routes == nil")
		} else {
			d := result2.Routes[0].Legs[0].Distance.Text

			if totalDistance, err := convertKm(d); err == nil {
				mapRes.Distance = fmt.Sprintf("%.4f", totalDistance)

			}
		}

		encoder := json.NewEncoder(w)
		encoder.Encode(&mapRes)

	}
}
