create table events (
    id serial,
    host_service_id integer,
    host_id integer,
    event_type varchar,
    service_name varchar(255),
    host_name varchar(255),
    message varchar(512),
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    primary key (id),
    foreign key (host_service_id) references host_services (id) on delete cascade on update cascade,
    foreign key (host_id) references hosts (id) on delete cascade on update cascade
);

create trigger set_timestamp
    before update on events
    for each row
    execute procedure trigger_set_timestamp();


