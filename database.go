package main

import (
	"fmt"
	"log"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// MongoDB stores the details of the DB connection.
type MongoDB struct {
	DatabaseURL    string
	DatabaseName   string
	CollectionName string
}

/*
Student represents the main persistent data structure.
It is of the form:
{
	"name": <value>, 	e.g. "Tom"
	"age": <value>		e.g. 21
	"studentid": <value>		e.c. "id0"
}
type Student struct {
	ID        bson.ObjectId `bson:"_id,omitempty"`
	Name      string        `json:"name"`
	Age       int           `json:"age"`
	StudentID string        `json:"studentid"`
}
*/

/*
Init initializes the mongo storage.
*/
func (db *MongoDB) Init(IDname string) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	index := mgo.Index{
		Key:        []string{IDname},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	err = session.DB(db.DatabaseName).C(db.CollectionName).EnsureIndex(index)
	if err != nil {
		panic(err)
	}
}

/*
AddTrack adds new track to the storage.
*/
func (db *MongoDB) AddTrack(s TrackData) error {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(db.DatabaseName).C(db.CollectionName).Insert(s)

	if err != nil {
		log.Printf("error in Insert(): %v", err.Error())
		return err
	}

	return nil
}

/*
Count returns the current count of the tracks in in-memory storage.
*/
func (db *MongoDB) Count() int {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// handle to "db"
	count, err := session.DB(db.DatabaseName).C(db.CollectionName).Count()
	if err != nil {
		fmt.Printf("error in Count(): %v", err.Error())
		return -1
	}

	return count
}

/*
GetTrack returns a track with a given ID or empty track struct.
*/
func (db *MongoDB) GetTrack(keyID string) (TrackData, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	track := TrackData{}
	allWasGood := true

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{trackDbIDName: keyID}).One(&track)
	if err != nil {
		allWasGood = false
	}

	return track, allWasGood
}

/*
GetAllTracks returns a slice with all the tracks.
*/
func (db *MongoDB) GetAllTracks() []TrackData {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var all []TrackData

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{}).All(&all)
	if err != nil {
		log.Print("db.GetAll error: ", err.Error())
		return []TrackData{}
	}

	return all
}

//--------Webhooks-----------

/*
AddHook adds new track to the storage.
*/
func (db *MongoDB) AddHook(s WebhookData) error {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(db.DatabaseName).C(db.CollectionName).Insert(s)

	if err != nil {
		log.Printf("error in Insert(): %v", err.Error())
		return err
	}

	return nil
}

/*
GetHook returns a track with a given ID or empty track struct.
*/
func (db *MongoDB) GetHook(keyID string) (WebhookData, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	track := WebhookData{}
	allWasGood := true

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{webhookDbIDName: keyID}).One(&track)
	if err != nil {
		allWasGood = false
	}

	return track, allWasGood
}

//GetAllHooks returns a slice with all the webhooks.
func (db *MongoDB) GetAllHooks() []WebhookData {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var all []WebhookData

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{}).All(&all)
	if err != nil {
		log.Print("db.GetAll error: ", err.Error())
		return []WebhookData{}
	}

	return all
}

//DeleteID an object from the db
func (db *MongoDB) DeleteID(field string, id int64) bool {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	allWasGood := true

	err = session.DB(db.DatabaseName).C(db.CollectionName).Remove(bson.M{field: id})
	if err != nil {
		log.Print("Failed to delete from database: ", err.Error())
		allWasGood = false
	}
	return allWasGood
}

//DeleteAll deletes all records of this database. plz no abuse
func (db *MongoDB) DeleteAll() bool {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	allWasGood := true

	_, err = session.DB(db.DatabaseName).C(db.CollectionName).RemoveAll(nil)
	if err != nil {
		log.Print("Failed to delete from database: ", err.Error())
		allWasGood = false
	}
	return allWasGood
}
