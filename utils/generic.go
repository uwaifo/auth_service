package utils

import (
	"app-auth/db"
	"context"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver/bson"
)

var AndEmptyString = ""
var AndTrue = true
var AndFalse = false

//when you have a slice of string
// you want to remove a specific value from an index
func RemoveIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

// when you want to find the index of the
func FindIndex(s []string, element string) int {
	for p, v := range s {
		if v == element {
			return p
		}
	}
	return -1
}

func DestroyRecord(collectionType, recordId string) (bool, error) {

	foundRecord := false
	if recordId == "" {
		return foundRecord, nil
	}

	filter := bson.D{{"id", recordId}}

	ok, err := db.TeamCollection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println(err)
		return foundRecord, nil
	}
	fmt.Println(ok)
	foundRecord = true
	return foundRecord, nil

}
