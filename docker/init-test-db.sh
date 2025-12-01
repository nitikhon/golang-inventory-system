#!/bin/bash
# filepath: ./docker/init-test-db.sh

set -e

echo "Creating test database..."

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
    CREATE DATABASE test_inventory_db;
    GRANT ALL PRIVILEGES ON DATABASE test_inventory_db TO $POSTGRES_USER;
EOSQL

echo "Test database 'test_inventory_db' created successfully!"