set -e
sudo cp -r ~/plgd/data/certs/ .
sudo chown -R $USER:$USER certs/
mv certs/external/coap-gateway.crt  certs/external/firmware-coap.crt
mv certs/external/coap-gateway.key  certs/external/firmware-coap.key