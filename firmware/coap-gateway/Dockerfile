FROM golang:1.15-alpine
ARG root_directory=$GOPATH/src/github.com/DavidMarquezF/mf-cloud/firmware/coap-gateway
RUN apk add --no-cache curl git build-base
WORKDIR $root_directory
COPY go.mod go.sum ./
RUN go mod download
COPY . .

ARG service=firmware-coap
RUN go build -ldflags "-linkmode external -extldflags -static" -o /go/bin/$service ./cmd/service

FROM ubuntu:20.04 as service
RUN apt update
COPY --from=build /go/bin/$service /usr/local/bin/$service
COPY run.sh /usr/local/bin/run.sh

# global
ENV FQDN="localhost"

# coap-gateway
ENV FIRMWARE_COAP_PORT=5684
ENV FIRMWARE_COAP_ADDRESS="0.0.0.0:$COAP_GATEWAY_PORT"
ENV FIRMWARE_COAP_DISABLE_VERIFY_CLIENTS=true
ENV FIRMWARE_COAP_BLOCKWISE_TRANSFER_SZX=1024
ENV FIRMWARE_COAP_LOG_MESSAGES=true


ENTRYPOINT ["/usr/local/bin/run.sh"]
