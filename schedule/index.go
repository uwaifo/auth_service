package schedule

import (
	"app-auth/db"
	"app-auth/types"
	"app-auth/utils"
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"fmt"

	"github.com/jasonlvhit/gocron"
	"github.com/mongodb/mongo-go-driver/bson"
)

func Starter() {

	cronInterval := os.Getenv("CRON_INTERVAL")
	number, err := strconv.ParseUint(cronInterval, 10, 64)
	if err != nil {
		log.Print(err)
	}

	gocron.Every(number).Minute().Do(runTask)

	go gocron.Start()
}
func runTask() {
	todo, err := GetTaksCount()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(todo, "Scheduled task(s)")
	tasks := GetTaks()
	fmt.Println(len(tasks), "Task(s) due to run")
	for _, v := range tasks {
		//fmt.Println("Scheduled", v.ItemType, v.Action, v.ItemId, "by", v.ExecuteOn)
		execItem, err := executor(v.ItemType, v.Action, v.ItemId)
		if err != nil {
			log.Print(err)
		}

		// TODO Consider that by checking on execItem, we will not delete the schedule if the item is not found
		if execItem {
			DestroySchedule(v.Id)
		}

	}
}

func executor(taskType, taskAction, taskItemId string) (bool, error) {

	successfullExecution := false

	if taskType == "ORG" && taskAction == "DestroyOrganisationRecord" {
		// call the orgs utility

		orgUtil := utils.OrganisationUtil{Id: taskItemId}

		destroyedOrganisation, err := orgUtil.DestroyOrganisationRecord()
		if !destroyedOrganisation && err != nil {
			return successfullExecution, err
		}
	} else if taskType == "TEAM" && taskAction == "DestroyTeamRecord" {

		//teamUtil := utils.TeamUtil{TeamId: taskItemId}

		destroyTeam, err := utils.DestroyRecord(taskType, taskItemId) // teamUtil.DestroyTeamRecord()

		if !destroyTeam && err != nil {
			fmt.Println("failed")

			return successfullExecution, err
		}
	}

	successfullExecution = true

	return successfullExecution, nil

}

func DestroySchedule(scheduleId string) {

	destroyScheduleUtility := utils.ScheduleUtil{
		Id: scheduleId,
	}

	destroyScheduleUtility.DestroyScheduleRecord()
	// fmt.Println("deleting schedule", scheduleId)

}

func GetTaks() []types.Schedule {

	// scheduleInterval := 10

	fromDate := time.Now()
	//toDate := time.Now().Add(5 * time.Minute)

	taskFilter := bson.M{
		"executeon": bson.M{
			//"$gte": fromDate,
			"$lt": fromDate,
		},
	}

	cur, err := db.AuthScheduleCollection.Find(context.Background(), taskFilter)
	if err != nil {
		log.Println(err)
		return nil
	}

	defer cur.Close(context.TODO())

	var tasks []types.Schedule

	for cur.Next(context.TODO()) {
		var task types.Schedule
		err := cur.Decode(&task)
		if err != nil {
			log.Print(err)
		}

		tasks = append(tasks, task)
	}

	return tasks

}

func GetTaksCount() (int, error) {

	taskFilter := bson.D{{"retrycount", 0}}

	cur, err := db.AuthScheduleCollection.Find(context.Background(), taskFilter)
	if err != nil {
		log.Println(err)
		return 0, nil
	}

	defer cur.Close(context.TODO())

	var tasks []types.Schedule

	for cur.Next(context.TODO()) {
		var task types.Schedule
		err := cur.Decode(&task)
		if err != nil {
			log.Print(err)
		}

		tasks = append(tasks, task)
	}

	return len(tasks), nil

}
