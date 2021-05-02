package utils

import (
	"app-auth/db"
	"app-auth/types"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/mongodb/mongo-go-driver/bson"
	uuid "github.com/satori/go.uuid"
)

type ScheduleUtil struct {
	Id string `json:"id"`
}

func (scheduleUtil ScheduleUtil) DestroyScheduleRecord() (bool, error) {
	foundRecord := false

	if scheduleUtil.Id == "" {
		return foundRecord, nil
	}

	scheduleFilter := bson.M{
		"id": scheduleUtil.Id,
	}

	ok, err := db.AuthScheduleCollection.DeleteOne(context.Background(), scheduleFilter)
	if err != nil {
		log.Println(err)
		return foundRecord, nil
	}

	fmt.Println(ok)
	foundRecord = true
	return foundRecord, nil

}

// ScheduleTask . . . .
func (scheduleUtil ScheduleUtil) ScheduleTask(taskType, taskAction, taskItemId string) (types.Schedule, error) {
	//Schedule Permanent Deletion

	deletionDays := os.Getenv("PERMANENT_DELETE_TIME_DAYS")
	deleteTime, err := strconv.ParseInt(deletionDays, 10, 64)
	if err != nil {
		log.Fatal(err)
		//return nil, err
	}

	scheduleId := uuid.NewV4().String()
	scheduleData := types.Schedule{
		Id:        scheduleId,
		ItemId:    taskItemId,
		ItemType:  taskType,
		Action:    taskAction,
		ExecuteOn: time.Now().Add(time.Duration(deleteTime) * time.Minute),
		CreatedAt: time.Now(),
	}
	response, err := db.AuthScheduleCollection.InsertOne(context.Background(), scheduleData)
	if err != nil {
		log.Println(err)
		//return err
	}
	fmt.Println(response)

	return scheduleData, nil
}
