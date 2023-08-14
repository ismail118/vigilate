package dbrepo

import (
	"github.com/ismail118/vigilate/internal/config"
	"github.com/ismail118/vigilate/internal/models"
	"github.com/ismail118/vigilate/internal/repository"
)

type testDBRepo struct {
	App *config.AppConfig
}

func NewTestDBRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}

// preferences
func (repo *testDBRepo) AllPreferences() ([]models.Preference, error) {
	var items []models.Preference
	return items, nil
}
func (repo *testDBRepo) SetSystemPref(name, value string) error {
	return nil
}
func (repo *testDBRepo) InsertOrUpdateSitePreferences(pm map[string]string) error {
	return nil
}
func (repo *testDBRepo) UpdateSystemPref(name, value string) error {
	return nil
}

// users and authentication
func (repo *testDBRepo) GetUserById(id int) (models.User, error) {
	var item models.User
	return item, nil
}
func (repo *testDBRepo) InsertUser(u models.User) (int, error) {
	return 0, nil
}
func (repo *testDBRepo) UpdateUser(u models.User) error {
	return nil
}
func (repo *testDBRepo) DeleteUser(id int) error {
	return nil
}
func (repo *testDBRepo) UpdatePassword(id int, newPassword string) error {
	return nil
}
func (repo *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 0, "", nil
}
func (repo *testDBRepo) AllUsers() ([]*models.User, error) {
	var items []*models.User
	return items, nil
}
func (repo *testDBRepo) InsertRememberMeToken(id int, token string) error {
	return nil
}
func (repo *testDBRepo) DeleteToken(token string) error {
	return nil
}
func (repo *testDBRepo) CheckForToken(id int, token string) bool {
	return true
}

// host
func (repo *testDBRepo) InsertHost(h models.Host) (int, error) {
	return 0, nil
}
func (repo *testDBRepo) GetHost(id int) (models.Host, error) {
	var item models.Host
	return item, nil
}
func (repo *testDBRepo) UpdateHost(h models.Host) error {
	return nil
}
func (repo *testDBRepo) GetListHosts() ([]models.Host, error) {
	var item []models.Host
	return item, nil
}

// host_services
func (repo *testDBRepo) UpdateHostServiceActive(hostID, serviceID, active int) error {
	return nil
}
func (repo *testDBRepo) GetAllServiceStatusCounts() (int, int, int, int, error) {
	return 0, 0, 0, 0, nil
}
func (repo *testDBRepo) GetServicesByStatus(status string) ([]models.Host, error) {
	var items []models.Host
	return items, nil
}
func (repo *testDBRepo) GetHostServiceByID(id int) (models.HostService, error) {
	var item models.HostService
	return item, nil
}
func (repo *testDBRepo) UpdateHostService(hs models.HostService) error {
	return nil
}
func (repo *testDBRepo) GetServicesToMonitor() ([]models.HostService, error) {
	var items []models.HostService
	return items, nil
}
func (repo *testDBRepo) GetHostServiceByHostIdServiceId(hostID, serviceID int) (models.HostService, error) {
	var item models.HostService
	return item, nil
}

// events
func (repo *testDBRepo) InsertEvent(e models.Event) (int, error) {
	return 0, nil
}
func (repo *testDBRepo) GetAllEvents() ([]models.Event, error) {
	var items []models.Event
	return items, nil
}