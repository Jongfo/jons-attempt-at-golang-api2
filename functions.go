package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	igc "github.com/marni/goigc"
)

//returns current timestamp
func timestampNow() int64 {
	return time.Since(baseTime).Nanoseconds() / 1000000
}

//errorHandler is a simple self made function to deal with bad requests
func errorHandler(w http.ResponseWriter, code int, mes string) {
	w.WriteHeader(code)
	http.Error(w, http.StatusText(code), code)
	fmt.Fprint(w, mes)
	log.Print(mes)
}

//trackDistance calculates rough distance
func trackDistance(t igc.Track) float64 {
	totalDistance := 0.0
	for j := 0; j < len(t.Points)-1; j++ {
		totalDistance += t.Points[j].Distance(t.Points[j+1])
	}
	return totalDistance
}

//copied code from stackoverflow. Could be improved on.
func diff(a, b time.Time) (year, month, day, hour, min, sec int) {
	if a.Location() != b.Location() {
		b = b.In(a.Location())
	}
	if a.After(b) {
		a, b = b, a
	}
	y1, M1, d1 := a.Date()
	y2, M2, d2 := b.Date()

	h1, m1, s1 := a.Clock()
	h2, m2, s2 := b.Clock()

	year = int(y2 - y1)
	month = int(M2 - M1)
	day = int(d2 - d1)
	hour = int(h2 - h1)
	min = int(m2 - m1)
	sec = int(s2 - s1)

	// Normalize negative values
	if sec < 0 {
		sec += 60
		min--
	}
	if min < 0 {
		min += 60
		hour--
	}
	if hour < 0 {
		hour += 24
		day--
	}
	if day < 0 {
		// days in month:
		t := time.Date(y1, M1, 32, 0, 0, 0, 0, time.UTC)
		day += 32 - t.Day()
		month--
	}
	if month < 0 {
		month += 12
		year--
	}

	return
}

func webhookPush() {
	//Go through all webhook
	for i := 0; i < len(webhookInfo); i++ {
		//check if number of new tracks >= webhook.mincap
		count := 0
		for j := 0; j < len(trackInfo); j++ {
			if trackInfo[j].timestamp > webhookInfo[i].Stop {
				count++
				if count >= webhookInfo[i].MinTriggerValue {
					fmt.Println("webhook POST")
					webhookInfo[i].Stop = trackInfo[j].timestamp

					var trids []string
					for k := 1; k <= count; k++ {
						trids = append(trids, trackInfo[j-count+k].UniqueID)
					}

					//Prepare message
					message := map[string]interface{}{
						//TODO: add prosessing time
						"content": fmt.Sprintf("Wow! New track was added!\nTimestamp: %d, New tracks: %+v", webhookInfo[i].Stop, trids),
					}
					count = 0

					bytesRepresentation, err := json.Marshal(message)
					if err != nil {
						log.Fatalln(err)
					}

					//send POST via url
					http.Post(webhookInfo[i].WebhookURL, "application/json", bytes.NewBuffer(bytesRepresentation))
				}
			}
		}
	}

}

//help function for clockticker(). No; this is not good practice.
//checks if z is less than x[index].timestamp and that index is not less than zero
func clockbool(z int64, x []TrackData, index int) bool {
	if index < 0 {
		return false
	}
	return z < x[index].timestamp
}

//make a ticker that runs every x minutes. POST to env var url
func clockticker() {
	//TODO: make data meaningful
	var lastTrackStamp int64
	if len(trackInfo) > 0 {
		lastTrackStamp = trackInfo[len(trackInfo)-1].timestamp
	} else {
		lastTrackStamp = 0
	}

	ticker := time.NewTicker(10 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				if len(trackInfo) > 0 {
					if lastTrackStamp < trackInfo[len(trackInfo)-1].timestamp {
						var trids []string
						//adds latest tracks to slice
						for i := len(trackInfo) - 1; clockbool(lastTrackStamp, trackInfo, i); i-- {
							trids = append(trids, trackInfo[i].UniqueID)
						}
						//invert slice for cronological order.
						for i, j := 0, len(trids)-1; i < j; i, j = i+1, j-1 {
							log.Printf("i: %d, j: %d", i, j)
							trids[i], trids[j] = trids[j], trids[i]
						}
						//Prepare message
						message := map[string]interface{}{
							"text": fmt.Sprintf("New tracks: %+v", trids),
						}
						bytesRepresentation, err := json.Marshal(message)
						if err != nil {
							log.Fatalln(err)
						}
						//send POST via url
						http.Post(os.Getenv("CLOUDCLOCKHOOK"), "application/json", bytes.NewBuffer(bytesRepresentation))
						lastTrackStamp = trackInfo[len(trackInfo)-len(trids)].timestamp
					}
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
}
