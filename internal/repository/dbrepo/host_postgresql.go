package dbrepo

import (
	"context"
	"github.com/ismail118/vigilate/internal/models"
	"log"
	"time"
)

func (m *postgresDBRepo) InsertHost(h models.Host) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	insert into hosts (host_name, canonical_name, url, ip, ipv6, location, os, active, created_at, updated_at) 
	values ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) returning id
`
	var newID int

	row := m.DB.QueryRowContext(ctx, query,
		h.HostName,
		h.CanonicalName,
		h.URL,
		h.IP,
		h.IPV6,
		h.Location,
		h.OS,
		h.Active,
		time.Now(),
		time.Now(),
	)
	err := row.Scan(&newID)
	if err != nil {
		return 0, err
	}

	// add new inactive host_services
	query = `
	insert into host_services 
	    (host_id, service_id, active, schedule_number, schedule_unit, status, last_check, created_at, updated_at)
	values 
	    ($1, $2, $3, $4, $5, $6, $7, $8, $9)
`

	_, err = m.DB.ExecContext(ctx, query,
		newID,
		1,
		0,
		3,
		"m",
		"ok",
		time.Now(),
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return newID, err
	}

	return newID, nil
}

func (m *postgresDBRepo) GetHost(id int) (models.Host, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select id, host_name, canonical_name, url, ip, ipv6, location, os, active, created_at, updated_at
	from hosts where id = $1
`
	var h models.Host

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&h.ID,
		&h.HostName,
		&h.CanonicalName,
		&h.URL,
		&h.IP,
		&h.IPV6,
		&h.Location,
		&h.OS,
		&h.Active,
		&h.CreatedAt,
		&h.UpdatedAt,
	)
	if err != nil {
		log.Println(err)
		return h, err
	}

	// get all services for host_service
	query = `
	select hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number, hs.schedule_unit, hs.status,
	       hs.last_check, hs.updated_at, hs.created_at, 
	       s.id, s.service_name, s.active, s.icon, s.created_at, s.updated_at
	from host_services hs
	left join services s on hs.service_id = s.id 
	where host_id = $1
`

	rows, err := m.DB.QueryContext(ctx, query, h.ID)
	if err != nil {
		return h, err
	}

	for rows.Next() {
		var hs models.HostService
		err = rows.Scan(
			&hs.ID,
			&hs.HostID,
			&hs.ServiceID,
			&hs.Active,
			&hs.ScheduleNumber,
			&hs.ScheduleUnit,
			&hs.Status,
			&hs.LastCheck,
			&hs.UpdatedAt,
			&hs.CreatedAt,
			&hs.Service.ID,
			&hs.Service.ServiceName,
			&hs.Service.Active,
			&hs.Service.Icon,
			&hs.Service.CreatedAt,
			&hs.Service.UpdatedAt,
		)
		if err != nil {
			return h, err
		}
		h.HostServices = append(h.HostServices, hs)
	}

	return h, nil
}

func (m *postgresDBRepo) UpdateHost(h models.Host) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	update hosts set host_name = $1, canonical_name = $2, url = $3, ip = $4,
	                 ipv6 = $5, location = $6, os = $7, active = $8, updated_at = $9
	where id = $10
`

	_, err := m.DB.ExecContext(ctx, query,
		h.HostName,
		h.CanonicalName,
		h.URL,
		h.IP,
		h.IPV6,
		h.Location,
		h.OS,
		h.Active,
		time.Now(),
		h.ID,
	)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func (m *postgresDBRepo) GetListHosts() ([]models.Host, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select id, host_name, canonical_name, url, ip, ipv6, location, os, active, created_at, updated_at
	from hosts
`
	items := make([]models.Host, 0)

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var h models.Host

		err := rows.Scan(
			&h.ID,
			&h.HostName,
			&h.CanonicalName,
			&h.URL,
			&h.IP,
			&h.IPV6,
			&h.Location,
			&h.OS,
			&h.Active,
			&h.CreatedAt,
			&h.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// get all services for host_service
		query = `
		select hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number, hs.schedule_unit, hs.status,
			   hs.last_check, hs.updated_at, hs.created_at, 
			   s.id, s.service_name, s.active, s.icon, s.created_at, s.updated_at
		from host_services hs
		left join services s on hs.service_id = s.id 
		where host_id = $1
`

		serviceRows, err := m.DB.QueryContext(ctx, query, h.ID)
		if err != nil {
			return nil, err
		}

		for serviceRows.Next() {
			var hs models.HostService
			err = serviceRows.Scan(
				&hs.ID,
				&hs.HostID,
				&hs.ServiceID,
				&hs.Active,
				&hs.ScheduleNumber,
				&hs.ScheduleUnit,
				&hs.Status,
				&hs.LastCheck,
				&hs.UpdatedAt,
				&hs.CreatedAt,
				&hs.Service.ID,
				&hs.Service.ServiceName,
				&hs.Service.Active,
				&hs.Service.Icon,
				&hs.Service.CreatedAt,
				&hs.Service.UpdatedAt,
			)
			if err != nil {
				return nil, err
			}
			h.HostServices = append(h.HostServices, hs)
		}
		serviceRows.Close()

		items = append(items, h)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return items, nil
}

// UpdateHostServiceActive update the active status of host_services
func (m *postgresDBRepo) UpdateHostServiceActive(hostID, serviceID, active int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	update host_services set active = $1 where host_id = $2 and service_id = $3
`

	_, err := m.DB.ExecContext(ctx, query, active, hostID, serviceID)
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) GetAllServiceStatusCounts() (int, int, int, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select
	(select count(id) from host_services where active = 1 and status = 'pending') as pending,
	(select count(id) from host_services where active = 1 and status = 'healthy') as healthy,
	(select count(id) from host_services where active = 1 and status = 'warning') as warning,
	(select count(id) from host_services where active = 1 and status = 'problem') as problem;
`
	var pending, healthy, warning, problem int

	row := m.DB.QueryRowContext(ctx, query)
	err := row.Scan(
		&pending,
		&healthy,
		&warning,
		&problem,
	)
	if err != nil {
		return 0, 0, 0, 0, err
	}

	return pending, healthy, warning, problem, nil
}

func (m *postgresDBRepo) GetServicesByStatus(status string) ([]models.Host, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number, hs.schedule_unit,
       hs.status, hs.last_check, hs.created_at, hs.updated_at,
       h.host_name, s.service_name
	from
		host_services hs
	left join hosts h on hs.host_id = h.id
	left join services s on hs.service_id = s.id
	where
		hs.status = $1
		and hs.active = 1
	order by 
	    host_name, service_name
	asc
`
	var items []models.Host

	rows, err := m.DB.QueryContext(ctx, query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var hs models.HostService
		var h models.Host
		err = rows.Scan(
			&hs.ID,
			&hs.HostID,
			&hs.ServiceID,
			&hs.Active,
			&hs.ScheduleNumber,
			&hs.ScheduleUnit,
			&hs.Status,
			&hs.LastCheck,
			&hs.CreatedAt,
			&hs.UpdatedAt,
			&h.HostName,
			&hs.Service.ServiceName,
		)
		if err != nil {
			return nil, err
		}

		h.HostServices = append(h.HostServices, hs)
		items = append(items, h)
	}

	return items, nil
}

func (m *postgresDBRepo) GetHostServiceByID(id int) (models.HostService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number, hs.schedule_unit, hs.status, 
	       hs.last_check, hs.updated_at, hs.created_at, 
	       s.id, s.service_name, s.active, s.icon, s.updated_at, s.created_at
	from 
	    host_services hs
	left join 
	        services s 
	    on hs.service_id = s.id
	where hs.id = $1
`
	var hs models.HostService
	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&hs.ID,
		&hs.HostID,
		&hs.ServiceID,
		&hs.Active,
		&hs.ScheduleNumber,
		&hs.ScheduleUnit,
		&hs.Status,
		&hs.LastCheck,
		&hs.UpdatedAt,
		&hs.CreatedAt,
		&hs.Service.ID,
		&hs.Service.ServiceName,
		&hs.Service.Active,
		&hs.Service.Icon,
		&hs.Service.UpdatedAt,
		&hs.Service.CreatedAt,
	)
	if err != nil {
		return hs, err
	}

	return hs, nil
}

func (m *postgresDBRepo) UpdateHostService(hs models.HostService) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	update host_services 
	set 
	    host_id = $1,
	    service_id = $2,
	    active = $3,
	    schedule_number = $4,
	    schedule_unit = $5,
	    status = $6,
	    last_check = $7,
	    updated_at = $8 
	where
	    id = $9
`

	_, err := m.DB.ExecContext(ctx, query,
		hs.HostID,
		hs.ServiceID,
		hs.Active,
		hs.ScheduleNumber,
		hs.ScheduleUnit,
		hs.Status,
		hs.LastCheck,
		hs.UpdatedAt,
		hs.ID,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *postgresDBRepo) GetServicesToMonitor() ([]models.HostService, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	query := `
	select hs.id, hs.host_id, hs.service_id, hs.active, hs.schedule_number, hs.schedule_unit, hs.status, 
	       hs.last_check, hs.updated_at, hs.created_at, 
	       s.id, s.service_name, s.active, s.icon, s.updated_at, s.created_at
	from 
	    host_services hs
	left join 
	        services s 
	    on hs.service_id = s.id
	where 
		hs.active = 1
`
	items := make([]models.HostService, 0)

	rows, err := m.DB.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var hs models.HostService

		err := rows.Scan(
			&hs.ID,
			&hs.HostID,
			&hs.ServiceID,
			&hs.Active,
			&hs.ScheduleNumber,
			&hs.ScheduleUnit,
			&hs.Status,
			&hs.LastCheck,
			&hs.UpdatedAt,
			&hs.CreatedAt,
			&hs.Service.ID,
			&hs.Service.ServiceName,
			&hs.Service.Active,
			&hs.Service.Icon,
			&hs.Service.UpdatedAt,
			&hs.Service.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		items = append(items, hs)
	}

	return items, nil
}