package database

/**
	Contains functions to perform various queries against database
*/

import (
	mgo "gopkg.in/mgo.v2"
	"time"
	"log"
)

var (
	Conn *mgo.Session
)

const (
	mgoHost = "72.47.237.205:27017"
	mgoDb = "dmvexchange"
	mgoUser = "dxb"
	mgoPwd = "9d9066cf90755496ec7ed3392e638fc111"
)

func InitMgo() {
	// Setup connection parameters
	dbDialInfo := &mgo.DialInfo{
		Addrs:    []string{mgoHost},
		Timeout:  60 * time.Second,
		Database: mgoDb,
		Username: mgoUser,
		Password: mgoPwd,
	}

	// attempt to connect
	session, err := mgo.DialWithInfo(dbDialInfo)
	if err != nil {
		log.Fatal(err)
	}
	session.SetMode(mgo.Strong, true)
	Conn = session
	return
}
