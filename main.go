package main

import (
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
)

var startTime time.Time
var baseTime time.Time
var db MongoDB
var idCap int

//trackInfo is a slice for all the tracks
var trackInfo []TrackData

func init() {
	//make a timestamp for uptime.
	startTime = time.Now()
	baseTime = time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC)
	idCap = 1
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
	//"/{route:route\\/?}"
	r.HandleFunc("/paragliding", redirAPI)
	r.HandleFunc("/paragliding/{api:api\\/?}", handlAPI)
	r.HandleFunc("/paragliding/api/track", handlAPItrack)
	r.HandleFunc("/paragliding/api/track/{ID}", handlAPItrackID)
	r.HandleFunc("/paragliding/api/track/{ID}/{field}", handlAPItrackIDfield)
	r.HandleFunc("/paragliding/api/ticker", handlAPIticker)
	r.HandleFunc("/paragliding/api/ticker/latest", handlAPItickerLatest)
	r.HandleFunc("/paragliding/api/ticker/{stamp}", handlAPItickerStamp)
	//	r.HandleFunc("/paragliding/api/webhook/new_track", handlAPIwebhookNT)
	//	r.HandleFunc("/paragliding/api/webhook/new_track/{WHID}", handlAPIwebhookNT)

	//serve our functionallity
	http.Handle("/", r)
	http.ListenAndServe(":"+port, nil)
}
