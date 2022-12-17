#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
	CREATE USER identity;
	CREATE DATABASE identity;
	GRANT ALL PRIVILEGES ON DATABASE identity TO identity;
	\connect identity;
  CREATE SCHEMA identity;
EOSQL