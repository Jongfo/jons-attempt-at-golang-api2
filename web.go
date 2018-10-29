package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	igc "github.com/marni/goigc"
)

//first self made function. Not relevant to task anymore
func handl404(w http.ResponseWriter, r *http.Request) {
	//sets header to 404
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "We found nothing exept this 404")
}

func redirAPI(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/paragliding/api", http.StatusSeeOther)
}

//writes meta information about our api
func handlAPI(w http.ResponseWriter, r *http.Request) {
	//make a timestamp, and compare it to startTime
	var tim time.Time
	tim = time.Now()
	y, mo, d, h, mi, s := diff(startTime, tim)
	//save the result as a string in the ISO8601 format.
	tim2 := fmt.Sprintf("P%dY%dM%dDT%dH%dM%dS",
		y,  //year
		mo, //month
		d,  //day
		h,  //hour
		mi, //min
		s)  //sec

	//make a json with metadata
	serv := Service{tim2, "Service for IGC tracks.", "v1"}
	js, err := json.Marshal(serv)
	if err != nil {
		str := fmt.Sprintf("Error Marshal: %s", err)
		errorHandler(w, http.StatusInternalServerError, str)
	} else {
		//forward json to requestee
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)
	}
}

func handlAPItrack(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		//store all IDs in a list. Allows empty json slice.
		ids := make([]string, 0)

		for i := 0; i < len(trackInfo); i++ {
			ids = append(ids, trackInfo[i].UniqueID)
		}
		//turn list into json
		js, err := json.Marshal(ids)
		if err != nil {
			str := fmt.Sprintf("Mershal error: %s", err)
			errorHandler(w, http.StatusInternalServerError, str)
		} else {
			//give json to customer
			w.Header().Set("Content-Type", "application/json")
			w.Write(js)
		}
	case "POST": //TODO: webhook check
		//make decoder for our POST body
		decoder := json.NewDecoder(r.Body)
		var url PostURL

		//decode the json we've recieved
		err1 := decoder.Decode(&url)
		if err1 != nil {
			str := fmt.Sprintf("Decode error: %s", err1)
			errorHandler(w, http.StatusInternalServerError, str)
			return //something went wrong
		}

		//get track information from provided url.
		track, err2 := igc.ParseLocation(url.URL)
		if err2 != nil {
			str := fmt.Sprintf("Problem reading the track: %s", err2)
			errorHandler(w, http.StatusInternalServerError, str)
			return //something went wrong
		}

		//put the ID of the track we found in a struct to be returned later.
		id := POSTid{track.UniqueID}
		//and turn it into a json
		js, err3 := json.Marshal(id)
		if err3 != nil {
			str := fmt.Sprintf("Marshal error: %s", err3)
			errorHandler(w, http.StatusInternalServerError, str)
			return //something went wrong
		}

		//check if we already have the track registered
		for i := 0; i < len(trackInfo); i++ {
			if trackInfo[i].UniqueID == track.UniqueID {
				str := fmt.Sprintf("Error: Already registered")
				errorHandler(w, http.StatusBadRequest, str)
				return //duplicate found
			}
		}

		//adds track to global slice
		//registeredTracks = append(registeredTracks, track)
		trackInfo = append(trackInfo, TrackData{track, url.URL, timestampNow()})

		//we did everything correctly, hopefully
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

		//webhook
		webhookPush()
	default:
		//unexpected request
		str := fmt.Sprintf("Sorry, only GET and POST methods are supported.")
		errorHandler(w, http.StatusBadRequest, str)
	}
}

//writes information about a track on a given ID
func handlAPItrackID(w http.ResponseWriter, r *http.Request) {
	//container for the http adress vaiable {ID}
	vars := mux.Vars(r)

	//goes throug and looks for matching ID
	for i := 0; i < len(trackInfo); i++ {
		if trackInfo[i].UniqueID == vars["ID"] {

			//calculates rough distance
			totalDistance := trackDistance(trackInfo[i].Track)

			//struct to be marshaled and sendt as json
			data := IDdata{
				trackInfo[i].Date,       //Date from File Header, H-record
				trackInfo[i].Pilot,      //Pilot name
				trackInfo[i].GliderType, //Glider type
				trackInfo[i].GliderID,   //Glider ID
				totalDistance,           //Calculated total track length
				trackInfo[i].url,
			}

			//make the json struct
			js, err := json.Marshal(data)
			if err != nil {
				str := fmt.Sprintf("Marshal error: %s", err)
				errorHandler(w, http.StatusInternalServerError, str)
			} else {
				//We successfully found data and made the json.
				w.Header().Set("Content-Type", "application/json")
				w.Write(js)
			}
			return //end function after finding track
		}
	}

	//in case we didn't find the track
	str := fmt.Sprintf("Error: Did not find track")
	errorHandler(w, http.StatusBadRequest, str)

}

//writes content of a specified ID and field
func handlAPItrackIDfield(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")

	//container for the http adress vaiables {ID} and {field}
	vars := mux.Vars(r)

	//Look for matching track ID
	for i := 0; i < len(trackInfo); i++ {
		if trackInfo[i].UniqueID == vars["ID"] {
			//Look for matching data name. If found; print data
			switch vars["field"] {
			case "pilot":
				fmt.Fprint(w, trackInfo[i].Pilot)
			case "glider":
				fmt.Fprint(w, trackInfo[i].GliderType)
			case "glider_id":
				fmt.Fprint(w, trackInfo[i].GliderID)
			case "track_length":
				fmt.Fprint(w, trackDistance(trackInfo[i].Track))
			case "H_date":
				fmt.Fprint(w, trackInfo[i].Date)
			case "track_src_url":
				fmt.Fprint(w, trackInfo[i].url)
			default:
				//last field does not match or not implemented yet.
				w.WriteHeader(http.StatusNotFound)
			}
			return
		}
	}
	//we did not find any matches
	w.WriteHeader(http.StatusNotFound)
}

//returns the latest timestamp of track
func handlAPItickerLatest(w http.ResponseWriter, r *http.Request) {
	if len(trackInfo) > 0 {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Fprint(w, trackInfo[len(trackInfo)-1].timestamp)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

//writes a json with timestamps and IDs
func handlAPIticker(w http.ResponseWriter, r *http.Request) {
	tNow := time.Now()
	//Check if no tracks have been added.
	if len(trackInfo) <= 0 {
		errorHandler(w, http.StatusSeeOther, "No tracks have been added yet.")
		return
	}
	var ids []string
	var tStop int64
	for i := 0; i < len(trackInfo) && i < idCap; i++ {
		ids = append(ids, trackInfo[i].UniqueID)
		tStop = trackInfo[i].timestamp
	}
	encoder := json.NewEncoder(w)

	jsun := TickerData{
		trackInfo[len(trackInfo)-1].timestamp, //latest
		trackInfo[0].timestamp,                //earliest
		tStop,                                 //last of ids
		ids,                                   //slice of UniqueID
		time.Since(tNow).Nanoseconds(),        //request time
	}

	w.Header().Set("Content-Type", "application/json")
	err := encoder.Encode(&jsun)
	if err != nil {
		errorHandler(w, http.StatusInternalServerError, "json encode error")
		return
	}
}

//writes a json with timestamps and IDs
func handlAPItickerStamp(w http.ResponseWriter, r *http.Request) {
	tNow := time.Now()
	//Check if no tracks have been added.
	if len(trackInfo) <= 0 {
		errorHandler(w, http.StatusSeeOther, "No tracks have been added yet.")
		return
	}
	//container for the http adress vaiable {stamp}
	vars := mux.Vars(r)
	index := -1
	//string to int64
	stamp, err := strconv.ParseInt(vars["stamp"], 10, 64)
	if err != nil {
		errorHandler(w, http.StatusBadRequest, "Expected number in adress.")
		return
	}

	//look for timestamp
	for i := 0; i < len(trackInfo); i++ {
		if trackInfo[i].timestamp == stamp {
			index = i
		}
	}
	//we might not have found any timestamp
	if index == -1 {
		errorHandler(w, http.StatusBadRequest, "Timestamp not found.")
		return
	}

	//prapare json data
	var ids []string
	var tStop int64
	for i := index + 1; i < len(trackInfo) && i < idCap+index+1; i++ {
		ids = append(ids, trackInfo[i].UniqueID)
		tStop = trackInfo[i].timestamp
	}
	encoder := json.NewEncoder(w)

	jsun := TickerData{
		trackInfo[len(trackInfo)-1].timestamp, //latest
		trackInfo[0].timestamp,                //earliest
		tStop,                                 //last of ids
		ids,                                   //slice of UniqueID
		time.Since(tNow).Nanoseconds(),        //request time
	}

	//write json data to body
	w.Header().Set("Content-Type", "application/json")
	err1 := encoder.Encode(&jsun)
	if err1 != nil {
		errorHandler(w, http.StatusInternalServerError, "json encode error")
		return
	}
}

//Webhook hendlers

func handlAPIwebhookNT(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		//Read request body
		decoder := json.NewDecoder(r.Body)
		var whData Webhookjson
		err := decoder.Decode(&whData)
		if err != nil {
			errorHandler(w, http.StatusInternalServerError, "Failed Decoding request")
			return
		}

		//prosses request
		if whData.MinTriggerValue <= 0 {
			fmt.Println("Min value not given") //debug
			whData.MinTriggerValue = 1
		}
		whID := timestampNow()
		webhookInfo = append(webhookInfo, WebhookData{whData, whID, whID})

		encoder := json.NewEncoder(w)

		//reply
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(201)
		encoder.Encode(fmt.Sprintf("%d", whID))
	default:
		errorHandler(w, http.StatusBadRequest, "Expected POST request")
	}
	return
}

func handlAPIwebhookID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//look for match
	if r.Method != "GET" && r.Method != "DELETE" {
		errorHandler(w, http.StatusBadRequest, "Expected GET or DELETE request")
		return
	}
	for i := 0; i < len(webhookInfo); i++ {
		if fmt.Sprintf("%d", webhookInfo[i].ID) == vars["WHID"] {
			//send response
			encoder := json.NewEncoder(w)
			w.Header().Set("Content-Type", "application/json")
			if encoder.Encode(webhookInfo[i].Webhookjson) != nil {
				errorHandler(w, http.StatusInternalServerError, "json encode error")
			}
			if r.Method == "DELETE" {
				//length := len(webhookInfo)
				var newSlice []WebhookData
				for j := 0; j < len(webhookInfo); j++ {
					if j != i {
						newSlice = append(newSlice, webhookInfo[j])
					}
				}
				log.Printf("We deleted webhook ID:%d", webhookInfo[i].ID)
				webhookInfo = newSlice
			}
			return
		}
	}
	//not found
	errorHandler(w, http.StatusNotFound, "We did not find any match")
}
