# jons-attempt-at-golang-api1
Assignment 2 for imt2681-2018

Project uses go modules. dependecies will be fetched on `go build`

Project uses optional environmental variables for some functionality.

##Environmental vriables
* `CLOUDCLOCKHOOK`: webhook link to slack webhook handler or slack compatible handler(Discord). unset or set to empty string to disable.
* `CLOUDCAP`: Default to 5 if unset. Changes the amount of track IDs returned for `/api/ticker`.


##Dependencies:
* github.com/golang/geo
* github.com/gorilla/mux
* github.com/kr/pretty
* github.com/marni/goigc
* gopkg.in/check.v1
* gopkg.in/mgo.v2