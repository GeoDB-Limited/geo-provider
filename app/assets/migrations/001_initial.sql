-- +migrate Up
create table devices
(
    id              bigserial,
    address         text      not null,
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
    PRIMARY KEY (id)
);

create table locations
(
    id        bigserial,
    address   text      not null,
    latitude  decimal   not null,
    longitude decimal   not null,
    altitude  decimal   not null,
    time      date      not null,
    timestamp timestamp not null,
    date      date      not null,
    PRIMARY KEY (id)
);

-- +migrate Down
drop table devices cascade;
drop table locations cascade;
