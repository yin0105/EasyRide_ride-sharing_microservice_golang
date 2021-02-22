package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type MapResp struct {
	Distance string `json:"distance"`
	ARoad    string `json:"a_road"`
}

type RosterResp struct {
	DriverCount string `json:"driver_count"`
	Rate        string `json:"rate"`
}

func handleStart(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	ss := strings.Split(uri, "/")
	for i, s := range ss {
		fmt.Println(i, s)
	}
	if len(ss) < 8 {
		http.Error(w, fmt.Sprintf("Wrong Parameters"), http.StatusBadRequest)
		return
	}
	origin := ss[4]
	destination := ss[5]
	startTimeHour := ss[6]
	startTimeMin := ss[7]

	var mapRes MapResp

	response, err := http.Get("http://localhost:8088/api/v1/mapping/origin=" + origin + "&destination=" + destination)
	if err != nil {
		http.Error(w, fmt.Sprintf("The HTTP request failed with error %s\n", err), http.StatusBadRequest)
		return
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		json.Unmarshal([]byte(data), &mapRes)
	}

	var rosterResp RosterResp

	response2, err2 := http.Get("http://localhost:8088/api/v1/driver/get_info/")
	if err2 != nil {
		http.Error(w, fmt.Sprintf("The HTTP request failed with error %s\n", err), http.StatusBadRequest)
		return
	} else {
		data, _ := ioutil.ReadAll(response2.Body)
		json.Unmarshal([]byte(data), &rosterResp)
	}

	var price string

	response3, err3 := http.Get("http://localhost:8088/api/v1/surge_pricing/" + mapRes.Distance + "/" + rosterResp.Rate + "/" + mapRes.ARoad + "/" + rosterResp.DriverCount + "/" + startTimeHour + "/" + startTimeMin)
	if err3 != nil {
		http.Error(w, fmt.Sprintf("The HTTP request failed with error %s\n", err), http.StatusBadRequest)
		return
	} else {
		data, _ := ioutil.ReadAll(response3.Body)
		json.Unmarshal([]byte(data), &price)
	}

	encoder := json.NewEncoder(w)
	encoder.Encode(price)
}
