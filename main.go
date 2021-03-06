package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const trackDbIDName = "trackid"
const webhookDbIDName = "hookid"

var startTime time.Time
var baseTime time.Time
var dbTrack MongoDB   //databade session object for tracks
var dbWebhook MongoDB //databade session object for webhooks
var idCap int         //Cap for trackID paging.        env variable
var clockhook string  //webhook link for clock tikker. env variable

//trackInfo is a slice for all the tracks
var trackInfo []TrackData
var webhookInfo []WebhookData

func init() {
	//make a timestamp for uptime.
	startTime = time.Now()
	baseTime = time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC)
	//environment variables
	idCap = 5
	i64, err := strconv.ParseInt(os.Getenv("CLOUDCAP"), 10, 64)
	if err == nil {
		idCap = int(i64)
	}
	log.Print("idcap: ", idCap)
	clockhook = os.Getenv("CLOUDCLOCKHOOK")
	log.Print("Hook url: ", clockhook)
}

func main() {
	//Database setup
	dbTrack = MongoDB{"mongodb://jongfo:the1kuk@ds135156.mlab.com:35156/joncloudtech", "joncloudtech", "trackcollection"}
	dbTrack.Init(trackDbIDName)
	dbWebhook = MongoDB{"mongodb://jongfo:the1kuk@ds135156.mlab.com:35156/joncloudtech", "joncloudtech", "webhookcollection"}
	dbWebhook.Init(webhookDbIDName)

	trackInfo = dbTrack.GetAllTracks()
	webhookInfo = dbWebhook.GetAllHooks()
	log.Print("Tracks fetched from database: ", len(trackInfo))
	log.Print("Webhooks fetched from database: ", len(webhookInfo))

	if os.Getenv("CLOUDCLOCKHOOK") != "" {
		clockticker()
	}

	//find our port
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	//all our paths using gorilla mux
	r := mux.NewRouter()
	r.HandleFunc("/", handl404)
	//"/{route:route\\/?}"{ID}
	r.HandleFunc("/{paragliding:paragliding\\/?}", redirAPI)
	r.HandleFunc("/paragliding/{api:api\\/?}", handlAPI)
	//track
	r.HandleFunc("/paragliding/api/{track:track\\/?}", handlAPItrack)
	r.HandleFunc("/paragliding/api/track/{ID}", handlAPItrackID)
	r.HandleFunc("/paragliding/api/track/{ID}/{field}", handlAPItrackIDfield)
	//ticker
	r.HandleFunc("/paragliding/api/ticker", handlAPIticker)
	r.HandleFunc("/paragliding/api/ticker/latest", handlAPItickerLatest)
	r.HandleFunc("/paragliding/api/ticker/{stamp}", handlAPItickerStamp)
	//webhook
	r.HandleFunc("/paragliding/api/webhook/new_track", handlAPIwebhookNT)
	r.HandleFunc("/paragliding/api/webhook/new_track/{WHID}", handlAPIwebhookID)
	//admin
	r.HandleFunc("/admin/api/{req}", handlAdmin)

	//serve our functionallity
	http.Handle("/", r)
	http.ListenAndServe(":"+port, nil)

	log.Print("End of service")
}
