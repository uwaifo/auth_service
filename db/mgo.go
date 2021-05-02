package db

import (
	"github.com/labstack/gommon/log"
	"os"

	"gopkg.in/mgo.v2"
)

func NewMongoSession() *mgo.Session {
	// if there's an environment variable available for the mongo url, use it.
	mongoHost := os.Getenv("MONGO_URL")

	session, err := mgo.Dial(mongoHost); if err != nil {
		log.Fatal(err)
	}

	return session
}