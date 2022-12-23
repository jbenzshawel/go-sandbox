#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER keycloak;
	ALTER USER keycloak WITH PASSWORD 'dockerdbpw';
	CREATE DATABASE keycloak;
	GRANT ALL PRIVILEGES ON DATABASE keycloak TO keycloak;
	\connect keycloak;
	GRANT ALL PRIVILEGES ON SCHEMA public TO keycloak;
	\connect postgres;
	CREATE USER identity;
	ALTER USER identity WITH PASSWORD 'dockerdbpw';
	CREATE DATABASE identity;
	GRANT ALL PRIVILEGES ON DATABASE identity TO identity;
	\connect identity;
    CREATE SCHEMA identity;
	\connect postgres;
EOSQL