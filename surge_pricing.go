package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

var distance, rate float64
var aRoad, driverCount, startTimeHour, startTimeMin int
var startTime string

func handleSurgePricing(w http.ResponseWriter, r *http.Request) {

	var err error
	uri := r.URL.Path
	ss := strings.Split(uri, "/")
	for i, s := range ss {
		fmt.Println(i, s)
	}
	if len(ss) < 10 {
		http.Error(w, fmt.Sprintf("Wrong Parameters"), http.StatusBadRequest)
		return
	}

	if distance, err = strconv.ParseFloat(ss[4], 32); err != nil {
		http.Error(w, fmt.Sprintf("Wrong Parameters"), http.StatusBadRequest)
		return
	}
	if rate, err = strconv.ParseFloat(ss[5], 32); err != nil {
		http.Error(w, fmt.Sprintf("Wrong Parameters"), http.StatusBadRequest)
		return
	}
	if aRoad, err = strconv.Atoi(ss[6]); err != nil {
		http.Error(w, fmt.Sprintf("Wrong Parameters"), http.StatusBadRequest)
		return
	}
	if driverCount, err = strconv.Atoi(ss[7]); err != nil {
		http.Error(w, fmt.Sprintf("Wrong Parameters"), http.StatusBadRequest)
		return
	}
	if startTimeHour, err = strconv.Atoi(ss[8]); err != nil {
		http.Error(w, fmt.Sprintf("Wrong Parameters"), http.StatusBadRequest)
		return
	}
	if startTimeMin, err = strconv.Atoi(ss[9]); err != nil {
		http.Error(w, fmt.Sprintf("Wrong Parameters"), http.StatusBadRequest)
		return
	}
	fmt.Sprintf("%f, %f, %i, %i, %i, %i", distance, rate, aRoad, driverCount, startTimeHour, startTimeMin)

	strTmp := string("0" + strconv.Itoa(startTimeHour))
	fmt.Println(strTmp)
	startTime = strTmp[len(strTmp)-2:]
	fmt.Println(startTime)
	strTmp = string("0" + strconv.Itoa(startTimeMin))
	startTime += ":" + strTmp[len(strTmp)-2:]

	fmt.Println(startTime)

	if price, ok := calcPrice(); ok {
		encoder := json.NewEncoder(w)
		encoder.Encode(fmt.Sprintf("%.4f", price))
	} else {
		http.Error(w, fmt.Sprintf("Could not calculate the surge price."), http.StatusBadRequest)
	}
}

func calcPrice() (float64, bool) {
	var price float64
	price = distance * rate
	if aRoad == 1 {
		price *= 2
	}
	if driverCount < 5 {
		price *= 2
	}
	if (startTime >= "23:00" && startTime <= "23:59") || (startTime >= "00:00" && startTime <= "06:00") {
		price *= 2
	}
	return price, true
}
