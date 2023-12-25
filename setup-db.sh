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

if [ -n "$POSTGRES_DB" ]; then
	create_enums $POSTGRES_DB;
	create_user_service_user $POSTGRES_USER_SERVICE_USER $POSTGRES_USER_SERVICE_PASS $POSTGRES_DB;
	create_ai_service_user $POSTGRES_AI_SERVICE_USER $POSTGRES_AI_SERVICE_PASS $POSTGRES_DB;
fi