package main

import (
	"database/sql"
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

func handleAddDriver(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	ss := strings.Split(uri, "/")
	for i, s := range ss {
		fmt.Println(i, s)
	}
	if len(ss) < 8 {
		http.Error(w, fmt.Sprintf("Wrong Parameters"), http.StatusBadRequest)
		return
	}
	user := ss[5]
	pwd := ss[6]
	rate := ss[7]

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/easy_ride")

	// if there is an error opening the connection, handle it
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect into MySQL database."), http.StatusBadRequest)
		panic(err.Error())
		return
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	// perform a db.Query insert
	insert, err := db.Query("INSERT INTO roster (user, pwd, rate) VALUES ( '" + user + "', '" + pwd + "', " + rate + " )")

	// if there is an error inserting, handle it
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to join a driver."), http.StatusBadRequest)
		panic(err.Error())
		return
	}
	// be careful deferring Queries if you are using transactions
	defer insert.Close()
	encoder := json.NewEncoder(w)
	encoder.Encode("OK")
}

func handleMapping(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Starting Mapping Microservice...")
	uri := r.URL.Path
	// response, err := http.Get("https://maps.googleapis.com/maps/api/directions/json?origin=37.75434337954133,%20-122.4837655029297&destination=137.750543040919084,%20122.41853417968751&key=AIzaSyDI57hkGB_K7Mtp4eFdYiy0mIw68z_1R1Y")
	response, err := http.Get("https://maps.googleapis.com/maps/api/directions/json?key=AIzaSyDI57hkGB_K7Mtp4eFdYiy0mIw68z_1R1Y&origin=37.75434337954133,%20-122.4837655029297&destination=37.750543040919084,%20-122.41853417968751")
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
				} else if unit == "ft" {
					s *= 0.0003048
				}
				fmt.Println(s) // 3.1415927410125732
			}
		}

	}
}
