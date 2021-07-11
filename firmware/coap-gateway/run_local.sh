export FQDN="localhost"

go build -ldflags "-linkmode external -extldflags -static" -o firmware-coap ./cmd/service


export NATS_PORT=10001
export RESOURCE_AGGREGATE_PORT=9083

export PLGD_DEV_FQDN=localhost #192.168.1.154


export HTTP_GATEWAY_PORT=9087

# coap-gateway
export FIRMWARE_COAP_PORT=5682
export FIRMWARE_COAP_ADDRESS="0.0.0.0:$FIRMWARE_COAP_PORT"
export FIRMWARE_COAP_DISABLE_VERIFY_CLIENTS=true
export FIRMWARE_COAP_BLOCKWISE_TRANSFER_SZX=1024
export FIRMWARE_COAP_LOG_MESSAGES=true
./run.sh
