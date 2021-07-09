#!/bin/bash                                                                     
DEVICE_ID="411590f1-905b-41e7-5325-a24ef97316c1"
MODULES='[{"id": 2},{"id":3}]'
PLATFORM=0
DEVICE_NAME="Test"

DATA='{"device_id":"'"$DEVICE_ID"'", "device_name":"'"$DEVICE_NAME"'", "platform":'"$PLATFORM"', "modules":'$MODULES'}'

echo $DATA

curl --header "Content-Type: application/json" \
  --request POST \
  --data "$DATA" \
  http://localhost:8091/createFirmware > test.elf
