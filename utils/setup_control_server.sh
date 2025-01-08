#!/usr/bin/env sh
# Copyright (c) 2025, Valentijn van de Beek <v.d.vandebeek@student.tudelft.nl>
# All rights reserved.
#
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

set -e
URL="localhost:4848"

# Set a preamble for each of the cURL requests. In particular:
#  - Content-Type -- Type of request data we are giving
#  - Origin -- Were the request originates from, at the moment we only accept localhost (needed for CORS)
#  - Type -- Allows a request to bypass session user checking
HEADERS=$(cat <<EOF
Content-Type: application/json
Origin: http://localhost:9090
Type: system
EOF
)

# Create the disk folder for the control server if it is not yet there
mkdir control_server/disks || true

# Setup an administrator user
curl -X POST "${URL}/user" -H "$HEADERS" -d \
     '{"username": "admin", "email": "EMAIL", "role": "administrator", "name": "Administrator" }'

# Create an image for the administrator user
curl -X POST "${URL}/user/admin/image" -H "$HEADERS" -d \
     '{"name": "Test image", "DiskCompressionStrategy": "none", "ImageFileType": "raw", "Type": "system", "username": "admin"}' | jq .UUID

# Setup a machine with the specified MAC address that is managed by BAAS
curl -X POST "${URL}/machine" -H "$HEADERS"  -d \
     "{\"name\": \"Test\", \"Architecture\": \"x86_64\", \"Managed\": true, \"MacAddress\": {\"Address\": \"$1\"}}"
