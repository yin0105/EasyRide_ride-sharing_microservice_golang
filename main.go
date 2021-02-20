package main

import (
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const port = 8087

// var verboseFlag bool

// func init() {
// 	flag.BoolVar(&verboseFlag, "verbose", false, "Enable verbose mode")
// 	flag.Parse()
// 	if verboseEnv := strings.ToLower(os.Getenv("MOVIE_VERBOSE")); verboseEnv == "true" || verboseFlag {
// 		log.SetLevel(log.DebugLevel)
// 	}
// 	formatter := &log.TextFormatter{
// 		FullTimestamp: true,
// 	}
// 	log.SetFormatter(formatter)
// }

func main() {
	http.HandleFunc("/api/v1/driver/add/", handleAddDriver)
	http.HandleFunc("/api/v1/driver/change/", handleChangeDriver)
	http.HandleFunc("/api/v1/driver/remove/", handleRemoveDriver)
	http.HandleFunc("/api/v1/drivers", handleDisplayDrivers)

	log.Printf("Starting movies microservice on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
