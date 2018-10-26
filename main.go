package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var startTime time.Time
var db MongoDB

//trackInfo is a slice for all the tracks
var trackInfo []TrackData

func init() {
	//make a timestamp for uptime.
	startTime = time.Now()
}

func main() {
	db = MongoDB{"mongodb://jongfo:the1kuk@ds135156.mlab.com:35156/joncloudtech", "joncloudtech", "trackcollection"}
	db.Init()

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
	r.HandleFunc("/paragliding/api/ticker", handlAPIticker)

	//serve our functionallity
	http.Handle("/", r)
	http.ListenAndServe(":"+port, nil)
}
