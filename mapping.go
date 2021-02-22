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
	Duration         Distance    `json:"duration"`
	EndLocation      Location    `json:"end_location"`
	HtmlInstructions string      `json:"html_instructions"`
	Polyline         interface{} `json:"polyline"`
	StartLocation    Location    `json:"start_location"`
	TravelMode       string      `json:"travel_mode"`
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
	ARoad    string `json:"a_road"`
}

type Address struct {
	Shop        string `json:"shop"`
	HouseNumber string `json:"house_number"`
	Road        string `json:"road"`
	City        string `json:"city"`
	County      string `json:"county"`
	State       string `json:"state"`
	Postcode    string `json:"postcode"`
	Country     string `json:"country"`
	CountryCode string `json:"country_code"`
}
type Nominatim struct {
	PlaceId     int64    `json:"place_id"`
	Licence     string   `json:"licence"`
	OsmType     string   `json:"osm_type"`
	OsmId       int64    `json:"osm_id"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	DisplayName string   `json:"display_name"`
	Address     Address  `json:"address"`
	Boundingbox []string `json:"boundingbox"`
}

func convertKm(d string) (float64, error) {
	d = strings.TrimSpace(d)
	unit := string(d[len(d)-2:])
	d = strings.TrimSpace(d[:len(d)-2])
	// fmt.Println("##" + d + "##")

	if s, err := strconv.ParseFloat(d, 64); err == nil {
		// fmt.Println(unit)
		// fmt.Println("%.4f", s)
		if unit == "mi" {
			s *= 1.60934
		} else if unit == "ft" {
			s *= 0.0003048
		}
		// fmt.Println("%.4f", s)
		return s, nil
	} else {
		fmt.Println("Error")
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
	// fmt.Println("https://maps.googleapis.com/maps/api/directions/json?key=AIzaSyDI57hkGB_K7Mtp4eFdYiy0mIw68z_1R1Y&" + uri[1])
	if err != nil {
		http.Error(w, fmt.Sprintf("The HTTP request failed with error %s\n", err), http.StatusBadRequest)
		return
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		// fmt.Println("data = " + string(data))
		var result map[string]interface{}
		json.Unmarshal([]byte(data), &result)

		// fmt.Println(result)

		var result2 HttpRes
		json.Unmarshal(data, &result2)
		if len(result2.Routes) == 0 {
			fmt.Println("Routes == nil")
		} else {
			d := result2.Routes[0].Legs[0].Distance.Text

			if totalDistance, err := convertKm(d); err == nil {
				mapRes.Distance = fmt.Sprintf("%.4f", totalDistance)
			}

			var total2 float64
			for _, step := range result2.Routes[0].Legs[0].Steps {
				fmt.Println("#################")

				if response2, err := http.Get("https://nominatim.openstreetmap.org/reverse?format=json&lat=" + fmt.Sprint(step.StartLocation.Lat) + "&lon=" + fmt.Sprint(step.StartLocation.Lng) + "&zoom=18&addressdetails=1"); err == nil {
					data2, _ := ioutil.ReadAll(response2.Body)
					var result2 Nominatim
					json.Unmarshal([]byte(data2), &result2)
					road := result2.Address.Road
					if len(road) > 0 && road[0] == 'A' {
						fmt.Println(result2.Address.Road)
					}
				}
				// fmt.Println(step.Distance.Text)
				if stepDistance, err := convertKm(step.Distance.Text); err == nil {
					fmt.Println(fmt.Sprintf("%.4f", stepDistance))
					total2 += stepDistance
				}
			}
			mapRes.ARoad = fmt.Sprintf("%.4f", total2)

		}

		encoder := json.NewEncoder(w)
		encoder.Encode(&mapRes)

	}
}
