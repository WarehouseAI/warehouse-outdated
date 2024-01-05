# /bin/sh
# multiple-psql-databases.sh

set -e
set -u

function create_ai_service_user() {
	local username=$1;
	local password=$2;
	local database=$3;

	echo "Creating user '$username'";

	psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$database" <<-EOSQL
		CREATE USER $username WITH PASSWORD '$password';
EOSQL
}

function create_user_service_user() {
	local username=$1;
	local password=$2;
	local database=$3;

	echo "Creating user '$username'"

	psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$database" <<-EOSQL
		CREATE USER $username WITH PASSWORD '$password';
EOSQL
}

function create_enums() {
	local database=$1;
	echo "Creating types in '$database'";

	psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$database" <<-EOSQL
			CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
	    CREATE TYPE authscheme AS ENUM ('Bearer', 'Basic','ApiKey');
			CREATE TYPE payloadtype AS ENUM ('JSON', 'FormData');
			CREATE TYPE iotype AS ENUM ('Image', 'Text', 'Audio');
			CREATE TYPE requestscheme AS ENUM ('POST', 'GET', 'PUT', 'UPDATE', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS');
			CREATE TYPE userrole AS ENUM ('Developer', 'Base');
EOSQL
}

function change_wal_level() {
	local database=$1;
	echo "Changing wal level in '$database' to logi";

	psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$database" <<-EOSQL
		ALTER SYSTEM SET wal_level = logical;
EOSQL
}

function create_publication() {
	local database=$1;
	local pub_name=$2;

	echo "Creating a publication on '$database' of all tables";
    
	psql -v ON_ERROR_STOP=1 -U "$POSTGRES_USER" -d "$database"  <<-EOSQL
		CREATE PUBLICATION $pub_name FOR All TABLES;
EOSQL
}

if [ -n "$POSTGRES_DB" ]; then
	create_enums $POSTGRES_DB;

	if [ $POSTGRES_DB = "ai_db" ]; then
		change_wal_level $POSTGRES_DB
		create_publication $POSTGRES_DB $AI_DB_PUB_NAME
	fi

	if [ $POSTGRES_DB = "users_db" ]; then
		change_wal_level $POSTGRES_DB
		create_publication $POSTGRES_DB $USERS_DB_PUB_NAME
	fi

	if [ $POSTGRES_DB = "stat_db" ]; then
		change_wal_level $POSTGRES_DB
	fi
fi