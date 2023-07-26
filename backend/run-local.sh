#!/bin/sh

# This file can be used to run the backend with a "real" home assistant instance
# for testing purposes (by default, the backend uses http://supervisor/core/api which
# is valid inside the addon container).
#
# You will need the API base URL of an instance and a long lived access token for it.

# Run:      ./run-local.sh <url> <token>
# E.g.:     ./run-local.sh "http://home.assistant:8123/api" "superSecretToken1337"

api_url="$1"
long_lived_token="$2"

SUPERVISOR_TOKEN=$long_lived_token go run . -hass-api-url=$api_url -cors-allow-all