create table host_services (
                          id serial,
                          host_id integer,
                          service_id integer,
                          active integer default 1,
                          schedule_number integer default 3,
                          schedule_unit varchar default 'm',
                          status varchar default 'pending',
                          last_check timestamp default '0001-01-01 00:00:01',
                          created_at timestamp default current_timestamp,
                          updated_at timestamp default current_timestamp,
                          primary key (id),
                          foreign key (host_id) references hosts (id) on delete cascade on update cascade ,
                          foreign key (service_id) references services (id) on delete cascade on update cascade
);

create trigger set_timestamp
    before update on host_services
    for each row
    execute procedure trigger_set_timestamp();


