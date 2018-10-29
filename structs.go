package main

import (
	"time"
)

/*
type TrackData struct {
	igc.Track        //"inherit from igc package"
	URL       string `json:"url"`
	Timestamp int64  `json:"timestamp"`
}*/

//TrackData contains all the relevant information about a track
type TrackData struct {
	UniqueID      string    `json:"uniqueid"`
	Date          time.Time `json:"date"`
	Pilot         string    `json:"pilot"`
	GliderType    string    `json:"glidertype"`
	GliderID      string    `json:"gliderid"`
	TotalDistance float64   `json:"totaldistance"`
	URL           string    `json:"url"`
	Timestamp     int64     `json:"timestamp"`
}

//WebhookData contains aditional information about a webhook
type WebhookData struct {
	Webhookjson
	ID   int64 `json:"timestamp"` //timestamp of when the webhook was made
	Stop int64 `json:"stop"`      //timestamp of the last gotten track
}

//-------json type structs------

//Webhookjson holds data of a webhook
type Webhookjson struct {
	WebhookURL      string `json:"webhookURL"`
	MinTriggerValue int    `json:"minTriggerValue"`
}

//Service contains data about our service
type Service struct {
	Uptime  string `json:"uptime"`
	Info    string `json:"info"`
	Version string `json:"version"`
}

//PostURL holds data on POST request.
type PostURL struct {
	URL string `json:"url"`
}

//POSTid contains ID based on PostURL
type POSTid struct {
	ID string `json:"id"`
}

//IDdata contains data on given track id
type IDdata struct {
	Hdate       time.Time `json:"H_date"`        //<date from File Header, H-record>,
	Pilot       string    `json:"pilot"`         //<pilot>,
	Glider      string    `json:"glider"`        //<glider>,
	GliderID    string    `json:"glider_id"`     //<glider_id>,
	TrackLength float64   `json:"track_length"`  //<calculated total track length>
	TrackURL    string    `json:"track_src_url"` //<the original URL used to upload the track>
}

//TickerData contains metadata about request and 5 track IDs.
type TickerData struct {
	TLatest    int64    `json:"t_latest"`   //<latest added timestamp>,
	TStart     int64    `json:"t_start"`    //<the first timestamp of the added track>, this will be the oldest track recorded
	TStop      int64    `json:"t_stop"`     //<the last timestamp of the added track>, this might equal to t_latest if there are no more tracks left
	TrackIDs   []string `json:"tracks"`     //[<id1>, <id2>, ...] cap at 5
	Processing int64    `json:"processing"` //<time in ns of how long it took to process the request>
}
