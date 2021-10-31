#!/bin/sh

until nc -z -v -w30 $MYSQL_HOST $MYSQL_PORT
do
  echo "Waiting for database connection..."
  # wait for 5 seconds before check again
  sleep 5
done

migrate \
  -path /migrations \
  -database "mysql://$MYSQL_USER:$MYSQL_PASSWORD@tcp($MYSQL_HOST:$MYSQL_PORT)/$MYSQL_DATABASE" \
  up
