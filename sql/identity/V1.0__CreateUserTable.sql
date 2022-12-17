CREATE TABLE IF NOT EXISTS identity.users (
     id serial,
     uuid UUID NOT NULL,
     first_name char(250) NOT NULL,
     last_name char(250) NOT NULL,
     email char(500) NOT NULL,
     PRIMARY KEY (id)
);