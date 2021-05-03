-- +migrate Up
create table devices
(
    address         text      not null unique,
    uuid            uuid      not null,
    os              text      not null,
    model           text      not null,
    locale          text      not null,
    apps            text      not null,
    version         text      not null,
    geocash_version text      not null,
    time            date      not null,
    timestamp       timestamp not null,
    date            date      not null,
    PRIMARY KEY (address)
);

create table locations
(
    address   text      not null unique,
    latitude  decimal   not null,
    longitude decimal   not null,
    altitude  decimal   not null,
    time      date      not null,
    timestamp timestamp not null,
    date      date      not null,
    PRIMARY KEY (address)
);

-- +migrate Down
drop table devices cascade;
drop table locations cascade;
