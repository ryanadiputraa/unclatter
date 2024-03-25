#!/bin/bash
# Extract command line arguments
server_port="$1"
server_fe_url="$2"

postgres_host="$3"
postgres_port="$4"
postgres_user="$5"
postgres_password="$6"
postgres_db="$7"

jwt_secret="$8"

google_redirect_url="$9"
google_client_id="$10"
google_client_secret="$11"
google_state="$12"


# Define the YAML content with placeholders replaced by command line arguments
YAML_CONTENT="
server:
  port: $server_port
  fe_url: $server_fe_url

postgres:
  host: $postgres_host
  port: $postgres_port
  user: $postgres_user
  password: $postgres_password
  db_name: $postgres_db
  ssl_mode: disable

jwt:
  secret: $jwt_secret

google_oauth:
  redirect_url: $google_redirect_url
  client_id: $google_client_id
  client_secret: $google_client_secret
  state: $google_state
"

# Write the YAML content to the file
echo "$YAML_CONTENT" > config/config.yml

echo "YAML config has been written to config/config.yml"
