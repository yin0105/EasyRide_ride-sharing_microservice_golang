package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type Driver struct {
	ID   int     `json:"id,string"`
	User string  `json:"user"`
	Rate float32 `json:"rate"`
}

type DriverInfo struct {
	DriverCount string `json:"driver_count"`
	Rate        string `json:"rate"`
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

func handleChangeDriver(w http.ResponseWriter, r *http.Request) {
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
	}

	results, err := db.Query("SELECT pwd FROM roster WHERE user='" + user + "'")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	if results.Next() {
		var curPwd string
		// for each row, scan the result into our tag composite object
		err = results.Scan(&curPwd)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		if curPwd != pwd {
			http.Error(w, fmt.Sprintf("Wrong password"), http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, fmt.Sprintf("No such user."), http.StatusBadRequest)
		return
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	// perform a db.Query update
	update, err := db.Query("UPDATE roster SET rate=" + rate + " where user='" + user + "' and pwd='" + pwd + "'")

	// if there is an error updating, handle it
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update rate."), http.StatusBadRequest)
		panic(err.Error())
		return
	}
	// be careful deferring Queries if you are using transactions
	defer update.Close()
	encoder := json.NewEncoder(w)
	encoder.Encode("OK")
}

func handleRemoveDriver(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	ss := strings.Split(uri, "/")
	for i, s := range ss {
		fmt.Println(i, s)
	}
	if len(ss) < 7 {
		http.Error(w, fmt.Sprintf("Wrong Parameters"), http.StatusBadRequest)
		return
	}
	user := ss[5]
	pwd := ss[6]

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/easy_ride")

	// if there is an error opening the connection, handle it
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect into MySQL database."), http.StatusBadRequest)
		panic(err.Error())
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	results, err := db.Query("SELECT pwd FROM roster WHERE user='" + user + "'")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	if results.Next() {
		var curPwd string
		// for each row, scan the result into our tag composite object
		err = results.Scan(&curPwd)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		if curPwd != pwd {
			http.Error(w, fmt.Sprintf("Wrong password"), http.StatusBadRequest)
			return
		}
	} else {
		http.Error(w, fmt.Sprintf("No such user."), http.StatusBadRequest)
		return
	}

	// perform a db.Query remove
	remove, err := db.Query("DELETE FROM roster where user='" + user + "' and pwd='" + pwd + "'")

	// if there is an error removing, handle it
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to leave a driver."), http.StatusBadRequest)
		panic(err.Error())
		return
	}
	// be careful deferring Queries if you are using transactions
	defer remove.Close()
	encoder := json.NewEncoder(w)
	encoder.Encode("OK")
}

func handleDisplayDrivers(w http.ResponseWriter, r *http.Request) {
	var drivers []Driver
	uri := r.URL.Path
	ss := strings.Split(uri, "/")
	for i, s := range ss {
		fmt.Println(i, s)
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/easy_ride")

	// if there is an error opening the connection, handle it
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect into MySQL database."), http.StatusBadRequest)
		panic(err.Error())
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	results, err := db.Query("SELECT id, user, rate FROM roster ORDER BY id")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for results.Next() {
		var driver Driver
		// for each row, scan the result into our tag composite object
		err = results.Scan(&driver.ID, &driver.User, &driver.Rate)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		drivers = append(drivers, driver)
	}
	encoder := json.NewEncoder(w)
	encoder.Encode(&drivers)
}

func handleGetInfo(w http.ResponseWriter, r *http.Request) {
	uri := r.URL.Path
	ss := strings.Split(uri, "/")
	for i, s := range ss {
		fmt.Println(i, s)
	}

	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/easy_ride")

	// if there is an error opening the connection, handle it
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to connect into MySQL database."), http.StatusBadRequest)
		panic(err.Error())
	}

	// defer the close till after the main function has finished executing
	defer db.Close()

	results, err := db.Query("SELECT COUNT(id), MIN(rate) FROM roster")
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	if results.Next() {
		var driverInfo DriverInfo
		// for each row, scan the result into our tag composite object
		err = results.Scan(&driverInfo.DriverCount, &driverInfo.Rate)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		// and then print out the tag's Name attribute
		encoder := json.NewEncoder(w)
		encoder.Encode(&driverInfo)
	}

}
