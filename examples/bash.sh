#!/bin/bash

STD_IN=$(</dev/stdin)
echo "PIPE STD_IN ${STD_IN}"
echo "hello error steam" > /dev/stderr
echo "hello out sream" > /dev/stderr
echo '{"status": 1, "reason": "ok", "payload": {}}'