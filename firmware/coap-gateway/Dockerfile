FROM golang:1.15-alpine as build
# Set necessary environmet variables needed for our image

RUN apk add --no-cache build-base

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Build the application
RUN go build -ldflags "-linkmode external -extldflags -static" -o main ./cmd/service

FROM ubuntu:20.04

RUN apt update
RUN apt install -y wget gnupg iproute2 systemctl openssl nginx ca-certificates
RUN wget https://github.com/mikefarah/yq/releases/download/v4.6.3/yq_linux_amd64 -O /usr/bin/yq && chmod +x /usr/bin/yq

# Copy binary from build to main folder
COPY --from=build /build/main firmware-coap
COPY run.sh /usr/local/bin/run.sh

# global
ENV FIRMWARE_COAP_FQDN="localhost"

# coap-gateway
ENV FIRMWARE_COAP_PORT=5682
ENV FIRMWARE_COAP_ADDRESS="0.0.0.0:$FIRMWARE_COAP_PORT"
ENV FIRMWARE_COAP_DISABLE_VERIFY_CLIENTS=true
ENV FIRMWARE_COAP_BLOCKWISE_TRANSFER_SZX=1024
ENV FIRMWARE_COAP_LOG_MESSAGES=true


CMD ["/usr/local/bin/run.sh"]
