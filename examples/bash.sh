#!/bin/bash

if [ "$1" == "schema" ]; then
    cat << EOF
{
    "name": "bash",
    "description": "Bash example",
    "inputs": {
       "\$schema": "http://json-schema.org/draft-07/schema#",
       "type": "object",
       "properties": {
          "user": {
             "type": "string"
          },
          "cmd": {
             "type": "string"
          },
          "requestId": {
             "type": "string"
          }
       },
       "required": [
          "user",
          "cmd",
          "requestId"
       ]
    },
    "outputs": {
       "$schema": "http://json-schema.org/draft-07/schema#",
       "type": "object",
       "properties": {
          "status": {
             "type": "integer"
          },
          "reason": {
             "type": "string"
          },
          "payload": {}
       },
       "required": [
          "status",
          "reason",
          "payload"
       ]
    }
 }
EOF
    exit 0
fi
STD_IN=$(</dev/stdin)
echo "PIPE STD_IN ${STD_IN}"
echo "hello error steam" > /dev/stderr
echo "hello out sream" > /dev/stderr
echo '{"status": 1, "reason": "ok", "payload": {}}'