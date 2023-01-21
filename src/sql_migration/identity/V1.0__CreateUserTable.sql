CREATE TABLE IF NOT EXISTS identity.users (
     user_id serial NOT NULL,
     user_uuid UUID NOT NULL,
     first_name varchar(250) NOT NULL,
     last_name varchar(250) NOT NULL,
     email varchar(500) NOT NULL,
     email_verified boolean NOT NULL,
     enabled boolean NOT NULL,
     created_at timestamp NOT NULL,
     last_updated_at timestamp NOT NULL,
     PRIMARY KEY (user_id)
);
