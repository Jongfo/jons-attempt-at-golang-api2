# jons-attempt-at-golang-api1
## Assignment 2 for imt2681-2018

Project uses go modules. dependecies will be fetched on `go build`

Project uses optional environmental variables for some functionality.

All webhooks use slack format. You can use discord webhooks by appending `/slack` at the end of the webhook url.


No trailing slashes!
## Paths
* `GET    /paragliding/api` 
Returns json with metat data about the API.

* `POST   /paragliding/api/track`
Takes a json with a url to a igc file and registers it. Retuns a json with a unique ID based on the igc file body.
POST Template:
```
{"url": "<url>"}
```
Response:
```
"id": "<id>"
```

* `POST   /paragliding/api/track`
returns json array with all registered IDs.
```
[<id1>, <id2>, ...]
```

* `GET    /paragliding/api/track/<id>`
Replace `<id>` with the text of a track ID to see meta information about this track. Displayed in application/json format.

* `GET    /paragliding/api/track/<id>/<field>`
Same as previous, exept that you only get the information about a spesific field in plain text.

* `GET    /paragliding/api/ticker/latest`
retuns the timestamp of the latest track in plain text.

* `GET    /paragliding/api/ticker`
returns a json struct with some timestamps and the 5(by default) first track IDs. 
```
{
"t_latest": <latest added timestamp>,
"t_start": <the first timestamp of the added track>, this will be the oldest track recorded
"t_stop": <the last timestamp of the added track>, this might equal to t_latest if there are no more tracks left
"tracks": [<id1>, <id2>, ...],
"processing": <time in ms of how long it took to process the request>
}

```

* `GET    /paragliding/api/ticker/<timestamp>`
Returns a json struct similar to above. Exept the returned tracks are after the given timestamp. Will return an error if no new tracks are found. 

* `POST   /paragliding/api/webhook/new_track`
Registers a webhook for notification about new tracks being added to the system. `webhookURL` expects a slack webhook. `minTriggerValue` is optional and defaults to 1.
Responds with code 201 if sccuessful.
Returns an ID for the webhook as plain text.
```
{
    "webhookURL": "url string",
    "minTriggerValue": "number"
}
```

* `GET    /paragliding/api/webhook/new_track/<webhook_id>`
Returns the contents of a saved webhook as it was registered in the above step.

* `DELETE /paragliding/api/webhook/new_track/<webhook_id>`
Does the same as above and deletes the registered webhook.


## Admin API:
* `GET    /admin/api/tracks_count`
Displays the amount of tracks saved locally and in the database.
* `DELETE /admin/api/tracks`
Deletes all the tracks saved locally and in the database
* `GET    /admin/api/webhooks_count`
Displays the amount of webhooks saved locally and in the database.
* `DELETE /admin/api/webhooks`
Deletes all the webhooks saved locally and in the database



## Heroku link
https://intense-bayou-18912.herokuapp.com/paragliding/api


## Environmental vriables
* `CLOUDCLOCKHOOK`: webhook link to slack webhook handler or slack compatible handler(Discord). unset or set to empty string to disable.
* `CLOUDCAP`: Default to 5 if unset. Changes the amount of track IDs returned for `/api/ticker`.


## Dependencies:
* github.com/golang/geo
* github.com/gorilla/mux
* github.com/kr/pretty
* github.com/marni/goigc
* gopkg.in/check.v1
* gopkg.in/mgo.v2