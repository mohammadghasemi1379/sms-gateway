#!/bin/sh

# Wait for MySQL to be ready
echo "Waiting for MySQL to be ready..."
while ! nc -z $DB_HOST $DB_PORT; do
  sleep 1
done
echo "MySQL is ready!"

# Wait for Redis to be ready
echo "Waiting for Redis to be ready..."
while ! nc -z $REDIS_HOST $REDIS_PORT; do
  sleep 1
done
echo "Redis is ready!"

# Start the application
echo "Starting SMS Gateway service..."
# Persist runtime environments
printenv >> /etc/environment
exec ./main
