#!/bin/bash                                                                     
DEVICE_ID="411590f1-905b-41e7-5325-a24ef97316c1"
TOKEN="Bearer eyJhbGciOiJFUzI1NiIsImtpZCI6IjUwMTE4MzQ5LTMwNzQtNTc0NS05ZTIwLWViYWZiZDY5Yzc1OSIsInR5cCI6IkpXVCJ9.eyJhdWQiOlsiaHR0cHM6Ly8xNzIuMjAuMS4yMTEvIl0KLCJjbGllbnRfaWQiOiJ0ZXN0IiwiZXhwIjoxNjI1NjE3Mjk1CiwiaWF0IjoxNjI1NjEzNjk1CiwiaXNzIjoiaHR0cHM6Ly8xNzIuMjAuMS4yMTEvIiwic2NvcGUiOlsib3BlbmlkIiwicjpkZXZpY2VpbmZvcm1hdGlvbjoqIiwicjpyZXNvdXJjZXM6KiIsInc6cmVzb3VyY2VzOioiLCJ3OnN1YnNjcmlwdGlvbnM6KiJdLCJzdWIiOiIxIn0.CVEwBu2yeOTxotV5uHWHpwqiJT5VBx3_7JtVFnZVrQJKkvfSvrcnZiuTZTQTgEekPfQDQWvaOkfyA3uL0AHdNw"
PLGD_CLOUD="https://172.20.1.211"
MODULES="[1,2,3]"

DATA='{"device_id":"'"$DEVICE_ID"'", "token":"'"$TOKEN"'", "plgd_cloud":"'"$PLGD_CLOUD"'", "modules":'$MODULES'}'

echo $DATA

curl --header "Content-Type: application/json" \
  --request POST \
  --data "$DATA" \
  http://localhost:8090/api/v1/fmw
