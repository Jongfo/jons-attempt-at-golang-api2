package main

import (
	"fmt"
	"log"
	"net/http"
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
