package main

import "strings"

type job struct {
	HostServiceID int
}

func (j job) Run() {

}

func startMonitoring() {
	if preferenceMap["monitoring_live"] == "1" {
		data := make(map[string]string)
		data["message"] = "starting"

		//TODO:

		// trigger a message to broadcast to all clients that app is starting to monitoring

		// get all of the services that we want to monitoring

		// rang throught the services

			// get the scheduler unit and number

			// create a job

			// save the id of the job so we can start/stop it

			// broadcast over web socket the fact that the services is scheduled
		
		// end range
	}
}