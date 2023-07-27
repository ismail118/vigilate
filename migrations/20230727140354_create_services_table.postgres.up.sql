create table services (
                          id integer,
                          service_name varchar(255),
                          active integer default 1,
                          icon varchar(255),
                          created_at timestamp default current_timestamp,
                          updated_at timestamp default current_timestamp,
                          primary key (id)
);

create trigger set_timestamp
    before update on services
    for each row
execute procedure trigger_set_timestamp();


