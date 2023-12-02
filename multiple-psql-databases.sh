# /bin/sh
# multiple-psql-databases.sh

set -e
set -u

function create_user_and_database() {
	local database=$1
	echo "Creating user and database '$database'"

	psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" <<-EOSQL
	    CREATE DATABASE $database;
	    GRANT ALL PRIVILEGES ON DATABASE $database TO $POSTGRES_USER;
EOSQL
}

function create_ai_enums() {
	local database=$1
	echo "Creating types in '$database'"

	psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$database" <<-EOSQL
			CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	    CREATE TYPE authscheme AS ENUM ('Bearer', 'Basic','ApiKey');
			CREATE TYPE payloadtype AS ENUM ('JSON', 'FormData');
			CREATE TYPE iotype AS ENUM ('Image', 'Text');
			CREATE TYPE requesttype AS ENUM ('POST', 'GET', 'PUT', 'UPDATE', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS');
EOSQL
}

function create_user_enums() {
	local database=$1
	echo "Creating types in '$database'"

	psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$database" <<-EOSQL
	    CREATE TYPE userrole AS ENUM ('Developer', 'Base');
EOSQL
}

if [ -n "$POSTGRES_MULTIPLE_DATABASES" ]; then
	echo "Multiple database creation requested: $POSTGRES_MULTIPLE_DATABASES"
	for db in $(echo $POSTGRES_MULTIPLE_DATABASES | tr ',' ' '); do
		create_user_and_database $db
		if [ $db = "ai" ]; then
			create_ai_enums $db
		fi

		if [ $db = "users" ]; then
			create_user_enums $db
		fi
	done
	echo "Multiple databases created"
fi