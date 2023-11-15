#!/bin/sh
# wait-for-postgres.sh

set -e

host="$1"
shift
cmd="$@"

until PGPASSWORD=$DB_PASS psql -h "$host" -U "$DB_USER" -d "$DB_NAME" -p $DB_PORT -c '\q'; do
  >&2 echo "Postgres in unavaliable - sleeping"
  sleep 1
done

>&2 echo "Postgres is up - executing command"
exec $cmd