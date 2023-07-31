create table hosts (
    id serial,
    host_name varchar(255),
    canonical_name varchar(255),
    url varchar(255),
    ip varchar(255),
    ipv6 varchar(255),
    location varchar(255),
    os varchar(255),
    active integer default 1,
    created_at timestamp default current_timestamp,
    updated_at timestamp default current_timestamp,
    primary key (id)
);

create trigger set_timestamp
    before update on hosts
    for each row
execute procedure trigger_set_timestamp();