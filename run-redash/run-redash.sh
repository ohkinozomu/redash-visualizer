#!/bin/bash

set -eu
docker-compose -f run-redash/docker-compose.yaml run --rm server create_db
docker-compose -f run-redash/docker-compose.yaml run --rm server /app/manage.py users create_root octocat@users.noreply.github.com root_user --password root_password
docker-compose -f run-redash/docker-compose.yaml run --rm server /app/manage.py ds new postgres --type pg --options '{"dbname": "postgres", "host": "127.0.0.1", "password": "postgres", "user": "password"}'
docker-compose -f run-redash/docker-compose.yaml up -d