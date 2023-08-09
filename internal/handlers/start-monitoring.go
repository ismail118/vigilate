package handlers

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type job struct {
	HostServiceID int
}

// Rn runs the schedulers
func (j job) Run() {
	Repo.ScheduledCheck(j.HostServiceID)
}

func (repo *DBRepo) StartMonitoring() {
	if app.PreferenceMap["monitoring_live"] == "1" {

		data := make(map[string]string)
		data["message"] = "Monitoring is starting..."
		err := app.WsClient.Trigger("public-channel", "app-starting", data)
		if err != nil {
			log.Println(err)
		}

		servicesToMonitor, err := repo.DB.GetServicesToMonitor()

		for _, hs := range servicesToMonitor {
			// get the scheduler unit and number
			var sch string
			if hs.ScheduleUnit == "d" {
				sch = fmt.Sprintf("@every %d%s", hs.ScheduleNumber*24, "h")
			} else {
				sch = fmt.Sprintf("@every %d%s", hs.ScheduleNumber, hs.ScheduleUnit)
			}

			// create a job
			job := job{
				HostServiceID: hs.ID,
			}

			scheduleID, err := app.Scheduler.AddJob(sch, job)
			if err != nil {
				log.Panicln(err)
			}

			// save the id of the job so we can start/stop it
			app.MonitorMap[hs.ID] = scheduleID

			// broadcast over web socket the fact that the services is scheduled
			payload := make(map[string]string)
			payload["message"] = "scheduling"
			payload["host_service_id"] = strconv.Itoa(hs.ID)
			yearOne := time.Date(0001, 11, 17, 20, 34, 58, 65138737, time.UTC)
			if app.Scheduler.Entry(app.MonitorMap[hs.ID]).Next.After(yearOne) {
				payload["next_run"] = app.Scheduler.Entry(app.MonitorMap[hs.ID]).Next.Format("2006-01-02 3:04:06 PM")
			} else {
				payload["next_run"] = "pending..."
			}
			payload["host"] = strconv.Itoa(hs.HostID)
			payload["service"] = hs.Service.ServiceName

			if hs.LastCheck.After(yearOne) {
				payload["last_run"] = hs.LastCheck.Format("2006-01-02 3:04:06 PM")
			} else {
				payload["last_run"] = "pending..."
			}

			payload["schedule"] = fmt.Sprintf("@every %d%s", hs.ScheduleNumber, hs.ScheduleUnit)
			err = app.WsClient.Trigger("public-channel", "next-run-event", payload)
			if err != nil {
				log.Println(err)
			}

			err = app.WsClient.Trigger("public-channel", "schedule-changed-event", payload)
			if err != nil {
				log.Println(err)
			}
		}
	}
}
