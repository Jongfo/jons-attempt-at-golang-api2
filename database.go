package main

import (
	"fmt"

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
func (db *MongoDB) Init() {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	index := mgo.Index{
		Key:        []string{"trackid"},
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
Add adds new track to the storage.
*/
func (db *MongoDB) Add(s TrackData) error {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	err = session.DB(db.DatabaseName).C(db.CollectionName).Insert(s)

	if err != nil {
		fmt.Printf("error in Insert(): %v", err.Error())
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
Get returns a track with a given ID or empty track struct.
*/
func (db *MongoDB) Get(keyID string) (TrackData, bool) {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	track := TrackData{}
	allWasGood := true

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{"trackid": keyID}).One(&track)
	if err != nil {
		allWasGood = false
	}

	return track, allWasGood
}

/*
GetAll returns a slice with all the tracks.
*/
func (db *MongoDB) GetAll() []TrackData {
	session, err := mgo.Dial(db.DatabaseURL)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	var all []TrackData

	err = session.DB(db.DatabaseName).C(db.CollectionName).Find(bson.M{}).All(&all)
	if err != nil {
		return []TrackData{}
	}

	return all
}
