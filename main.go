package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var startTime time.Time

//locally stored tracks. Depricated
//var registeredTracks []igc.Track

//trackInfo is a slice for all the tracks
var trackInfo []TrackData

func init() {
	//make a timestamp for uptime.
	startTime = time.Now()
}

func main() {
	//find our port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//all our paths using gorilla mux
	r := mux.NewRouter()
	r.HandleFunc("/", handl404)
	r.HandleFunc("/paragliding", redirAPI)
	r.HandleFunc("/paragliding/api", handlAPI)
	r.HandleFunc("/paragliding/api/track", handlAPItrack)
	r.HandleFunc("/paragliding/api/track/{ID}", handlAPItrackID)
	r.HandleFunc("/paragliding/api/track/{ID}/{field}", handlAPItrackIDfield)

	//serve our functionallity
	http.Handle("/", r)
	http.ListenAndServe(":"+port, nil)
}
