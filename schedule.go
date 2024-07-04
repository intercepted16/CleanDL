package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-co-op/gocron/v2"
)

func ScheduleDailyTask() {
	// create a scheduler
	s, err := gocron.NewScheduler()
	if err != nil {
		// handle error
		os.Exit(1)
	}

	// add a job to the scheduler
	_, err = s.NewJob(
		gocron.DailyJob(
			1, gocron.NewAtTimes(gocron.NewAtTime(0, 0, 0)),
		),
		gocron.NewTask(
			func() {
				downloadsFolder := getDownloadsFolder()
				patterns := getSettings(patternsPath)
				processFiles(patterns, downloadsFolder)
			},
		),
	)
	if err != nil {
		fmt.Println(err)
	}
	// start the scheduler
	s.Start()

	// block until you are ready to shut down
	<-time.After(time.Minute)

	// when you're done, shut it down
	err = s.Shutdown()
	if err != nil {
		os.Exit(1)
		// handle error
	}
}
