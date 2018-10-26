package main

import (
	"time"

	igc "github.com/marni/goigc"
)

//TrackData contains all the relevant information about a track
type TrackData struct {
	igc.Track //"inherit from igc package"
	url       string
}

//-------json type structs------

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
