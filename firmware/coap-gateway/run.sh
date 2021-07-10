
#!/usr/bin/env bash
#set -e



export CERTIFICATES_PATH="./certs"
export LOGS_PATH="./log"

export INTERNAL_CERT_DIR_PATH="$CERTIFICATES_PATH/internal"
export GRPC_INTERNAL_CERT_NAME="endpoint.crt"
export GRPC_INTERNAL_CERT_KEY_NAME="endpoint.key"

# ENDPOINTS
export NATS_HOST="$PLGD_DEV_FQDN:$NATS_PORT"
export NATS_URL="nats://${NATS_HOST}" 

# ROOT CERTS
export CA_POOL_DIR="$CERTIFICATES_PATH"
export CA_POOL_NAME_PREFIX="root_ca"
export CA_POOL_CERT_PATH="$CA_POOL_DIR/$CA_POOL_NAME_PREFIX.crt"
export CA_POOL_CERT_KEY_PATH="$CA_POOL_DIR/$CA_POOL_NAME_PREFIX.key"


#LISTEN CERTS
export LISTEN_TYPE="file"
export LISTEN_FILE_CA_POOL="$CA_POOL_CERT_PATH"


export EXTERNAL_CERT_DIR_PATH="$CERTIFICATES_PATH/external"

export FIRMWARE_COAP_FILE_CERT_NAME="firmware-coap.crt"
export FIRMWARE_COAP_FILE_CERT_KEY_NAME="firmware-coap.key"


mkdir -p $LOGS_PATH


echo "starting coap-gateway-secure"
#need to export variable to make it available to subprocess
# Blockwise has to be enabled in order to get the firmware
export ADDRESS=${FIRMWARE_COAP_ADDRESS}
export LOG_MESSAGES=${FIRMWARE_COAP_LOG_MESSAGES}
export EXTERNAL_PORT=${FIRMWARE_COAP_PORT}
export FQDN=${FIRMWARE_COAP_FQDN}
export LISTEN_FILE_CERT_DIR_PATH=${EXTERNAL_CERT_DIR_PATH}
export LISTEN_FILE_CERT_NAME=${FIRMWARE_COAP_FILE_CERT_NAME}
export LISTEN_FILE_CERT_KEY_NAME=${FIRMWARE_COAP_FILE_CERT_KEY_NAME}
export LISTEN_FILE_DISABLE_VERIFY_CLIENT_CERTIFICATE=${FIRMWARE_COAP_DISABLE_VERIFY_CLIENTS}
export DISABLE_BLOCKWISE_TRANSFER=false 
export BLOCKWISE_TRANSFER_SZX=${FIRMWARE_COAP_BLOCKWISE_TRANSFER_SZX}
export DISABLE_PEER_TCP_SIGNAL_MESSAGE_CSMS=${FIRMWARE_COAP_DISABLE_PEER_TCP_SIGNAL_MESSAGE_CSMS} 
./firmware-coap >$LOGS_PATH/firmware-coap.log 2>&1 
status=$?
if [ $status -ne 0 ]; then
  echo "Failed to start coap-gateway: $status"
  sync
  cat $LOGS_PATH/firmware-coap.log
  exit $status
fi




