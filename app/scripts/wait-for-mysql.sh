#!/bin/sh

set -e

until mysql -h db -u "$DB_USERNAME" -p"$DB_PASSWORD" "$DB_NAME" -e "\q"; do
    echo >&2 "MySQL is unavailable - sleeping"
    sleep 1
done

echo >&2 "MySQL is up - executing command"
exec "$@"
