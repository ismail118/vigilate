package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/ismail118/vigilate/internal/models"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	HTTP           = 1
	HTTPS          = 2
	SSLCertificate = 3
	StatusProblem  = "problem"
	StatusHealthy  = "healthy"
)

type jsonResp struct {
	OK            bool      `json:"ok"`
	Message       string    `json:"message"`
	ServiceID     int       `json:"service_id"`
	HostServiceID int       `json:"host_service_id"`
	HostID        int       `json:"host_id"`
	OldStatus     string    `json:"old_status"`
	NewStatus     string    `json:"new_status"`
	LastCheck     time.Time `json:"last_check"`
}

func (repo *DBRepo) ScheduledCheck(hostServiceID int) {
	log.Println("----------- running check for", hostServiceID)

	hs, err := repo.DB.GetHostServiceByID(hostServiceID)
	if err != nil {
		log.Println(err)
		return
	}

	h, err := repo.DB.GetHost(hs.HostID)
	if err != nil {
		log.Println(err)
		return
	}

	newStatus, msg := repo.testServiceForHost(h, hs)

	if newStatus != hs.Status {
		repo.updateHostServiceStatusCount(h, hs, newStatus, msg)
	}

	log.Println(newStatus, "=", msg)
}

func (repo *DBRepo) updateHostServiceStatusCount(h models.Host, hs models.HostService, newStatus, msg string) {
	// update host_service field status and last check 
	hs.Status = newStatus
	hs.LastCheck = time.Now()
	err := repo.DB.UpdateHostService(hs)
	if err != nil {
		log.Println(err)
		return
	}

	pending, healthy, warning, problem, err := repo.DB.GetAllServiceStatusCounts()
	if err != nil {
		log.Println(err)
		return
	}
	data := make(map[string]string)
	data["pending_count"] = strconv.Itoa(pending)
	data["healthy_count"] = strconv.Itoa(healthy)
	data["warning_count"] = strconv.Itoa(warning)
	data["problem_count"] = strconv.Itoa(problem)
	repo.broadcastMessage("public-channel", "host-service-count-changed", data)
}

func (repo *DBRepo) broadcastMessage(channel, event string, data map[string]string) {
	err := repo.App.WsClient.Trigger(channel, event, data)
	if err != nil {
		log.Println(err)
	}
}

func (repo *DBRepo) TestCheck(w http.ResponseWriter, r *http.Request) {
	hostServiceId, _ := strconv.Atoi(chi.URLParam(r, "id"))
	oldStatus := chi.URLParam(r, "oldStatus")
	ok := true

	// get host_services
	hs, err := repo.DB.GetHostServiceByID(hostServiceId)
	if err != nil {
		log.Println(err)
		ok = false
	}

	// get host
	h, err := repo.DB.GetHost(hs.HostID)
	if err != nil {
		log.Println(err)
		ok = false
	}

	// test the service
	newStatus, msg := repo.testServiceForHost(h, hs)

	// update the host_services status (if changed) and last check
	hs.Status = newStatus
	hs.LastCheck = time.Now()
	hs.UpdatedAt = time.Now()

	// broadcast service status changed event
	err = repo.DB.UpdateHostService(hs)
	if err != nil {
		log.Println(err)
		ok = false
	}

	// create json
	var resp jsonResp
	if ok {
		resp = jsonResp{
			OK:            true,
			Message:       msg,
			ServiceID:     hs.ServiceID,
			HostServiceID: hs.ID,
			HostID:        hs.HostID,
			OldStatus:     oldStatus,
			NewStatus:     newStatus,
			LastCheck:     time.Now(),
		}
	} else {
		resp = jsonResp{
			OK:      true,
			Message: "Something went wrong",
		}
	}

	// send json to client
	out, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (repo *DBRepo) SetSystemPref(w http.ResponseWriter, r *http.Request) {
	resp := jsonResp{
		OK: true,
		Message: "Success",
	}

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		resp.OK = false
		resp.Message = err.Error()
	}

	prefName := r.PostForm.Get("pref_name")
	prefValue := r.PostForm.Get("pref_value")
	err = repo.DB.UpdateSystemPref(prefName, prefValue)
	if err != nil {
		log.Println(err)
		resp.OK = false
		resp.Message = err.Error()
	}

	repo.App.PreferenceMap["monitoring_live"] = prefValue

	out, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}


// ToggleMonitoring turns monitoring on and off
func (repo *DBRepo) ToggleMonitoring(w http.ResponseWriter, r *http.Request) {
	resp := jsonResp{
		OK: true,
		Message: "Success",
	}

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		resp.OK = false
		resp.Message = err.Error()
	}

	enabled := r.PostForm.Get("enabled")

	if enabled == "1" {
		// start monitoring
		log.Println("turn monitoring on")
		repo.App.PreferenceMap["monitoring_live"] = enabled
		repo.StartMonitoring()
		// start the scheduler
		repo.App.Scheduler.Start()
	} else {
		// stop montirong
		log.Println("turn monitoring off")

		repo.App.PreferenceMap["monitoring_live"] = enabled

		// remove all items in map from scheduler and delete from map
		for k, v := range repo.App.MonitorMap {
			repo.App.Scheduler.Remove(v)
			delete(repo.App.MonitorMap, k)
		}

		// delete all entries from scheduler, to be sure
		for _, x := range repo.App.Scheduler.Entries() {
			repo.App.Scheduler.Remove(x.ID)
		}

		// emptyh the map
		repo.App.Scheduler.Stop()

		data := make(map[string]string)
		data["message"] = "Monitoring is off!"
		err := app.WsClient.Trigger("public-channel", "app-stopping", data)
		if err != nil {
			log.Println(err)
		}
	}

	out, _ := json.Marshal(resp)
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

func (repo *DBRepo) testServiceForHost(h models.Host, hs models.HostService) (string, string) {
	var msg, newStatus string

	switch hs.ServiceID {
	case HTTP:
		msg, newStatus = testHTTPForHost(h.URL)
	}

	// broadcast to clients
	if hs.Status != newStatus {
		data := make(map[string]string)
		data["host_id"] = strconv.Itoa(hs.HostID)
		data["host_service_id"] = strconv.Itoa(hs.ID)
		data["host_name"] = h.HostName
		data["service_name"] = hs.Service.ServiceName
		data["icon"] = hs.Service.Icon
		data["status"] = newStatus
		data["active"] = strconv.Itoa(hs.Active)
		data["message"] = fmt.Sprintf("%s on %s reports %s", hs.Service.ServiceName, h.HostName, newStatus)
		data["last_check"] = time.Now().Format("2006-01-02 3:04:06 PM")
		repo.broadcastMessage("public-channel", "host-service-status-changed", data)
	}

	//TODO: if appropriate send email or sms message
	return newStatus, msg
}

func testHTTPForHost(url string) (string, string) {
	if strings.HasSuffix(url, "/") {
		url = strings.TrimSuffix(url, "/")
	}

	url = strings.Replace(url, "https://", "http://", -1)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Sprintf("%s - %s", url, "error connecting"), StatusProblem
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("%s - %s", url, resp.Status), StatusProblem
	}

	return fmt.Sprintf("%s - %s", url, resp.Status), StatusHealthy
}