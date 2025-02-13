CREATE TABLE urltable(
    id serial unique not null,
    full_url varchar(255),
    short_url varchar(255)
);