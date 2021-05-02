package main

//import (
//	"app-auth/config"
//	"app-auth/db"
//	"app-auth/utils"
//	"fmt"
//	"github.com/labstack/gommon/log"
//	"github.com/mongodb/mongo-go-driver/bson"
//)
//
//var collectionName, _ = utils.CreateCollection(config.ConfirmationDbName)
//var collection, _ = utils.CreateCollection("users")
//
//
//func main() {
//	//http://localhost:5000/confirm/a283573e-223f-4e3a-9d5b-cf4315bc0624
//	filter := bson.M{"confirmid": "a283573e-223f-4e3a-9d5b-cf4315bc0624"}
//	holder := config.ConfirmationData{}
//
//	err := db.FindOne(collectionName, utils.CreateTimeoutContext(5), filter).Decode(&holder); if err != nil {
//		fmt.Println(err)
//	}
//
//	fmt.Println(utils.Stringify(holder))
//
//	filterUser := bson.M{"email": holder.Email, "id": holder.Userid, "provider": "local"}
//	resultContainer := config.UserData{}
//
//	err = db.FindOne(collection, utils.CreateTimeoutContext(5), filterUser).Decode(&resultContainer); if err != nil {
//		// user is not found return something
//		log.Error(err)
//	}
//
//	fmt.Println(utils.Stringify(resultContainer))
//
//	resultContainer.Verified = false
//
//	userData := utils.UserBsonM(resultContainer)
//
//	fmt.Println(utils.Stringify(userData))
//
//	err = db.FindOneAndReplace(collection, utils.CreateTimeoutContext(5), filterUser, userData).Decode(&resultContainer); if err != nil {
//		log.Error(err)
//	}
//
//	log.Info(utils.Stringify(resultContainer))
//}
