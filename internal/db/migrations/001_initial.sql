-- +migrate Up

CREATE TABLE locations (
                        address    CHARACTER(42) not null,
                    latitude       bigint not null,
                        longitude       bigint not null,
                        altitude       bigint not null,
                        time TIMESTAMP
);

create index locations_address_index on locations(address);
create index locations_latitude_index on locations(latitude);
create index locations_longitude_index on locations(longitude);
create index locations_altitude_index on locations(altitude);

-- +migrate Down

drop table if exists locations;