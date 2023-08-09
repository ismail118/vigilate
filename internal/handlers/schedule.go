package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"

	"github.com/CloudyKit/jet/v6"
	"github.com/ismail118/vigilate/internal/helpers"
	"github.com/ismail118/vigilate/internal/models"
)

type ByHost []models.Schedule

// Len is used to sort by host
func (b ByHost) Len() int {
	return len(b)
}

// Less	 is used to sort by host
func (b ByHost) Less(i, j int) bool {
	return b[i].Host < b[j].Host
}

// Swap is used to sort by host
func (b ByHost) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}

// ListEntries lists schedule entries
func (repo *DBRepo) ListEntries(w http.ResponseWriter, r *http.Request) {
	var items []models.Schedule

	for k, v := range repo.App.MonitorMap {
		var item models.Schedule

		item.ID = k
		item.EntryID = v
		item.Entry = repo.App.Scheduler.Entry(v)
		hs, err := repo.DB.GetHostServiceByID(k)
		if err != nil {
			log.Println(err)
			return
		}
		item.ScheduleText = fmt.Sprintf("@every %d%s", hs.ScheduleNumber, hs.ScheduleUnit)
		item.LastRunFromHS = hs.LastCheck
		item.Host = strconv.Itoa(hs.HostID)
		item.Service = hs.Service.ServiceName
		items = append(items, item)
	}

	// sort slice by host
	sort.Sort(ByHost(items))

	vars := make(jet.VarMap)
	vars.Set("items", items)

	err := helpers.RenderPage(w, r, "schedule", vars, nil)
	if err != nil {
		printTemplateError(w, err)
	}
}
