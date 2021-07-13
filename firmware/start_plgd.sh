# Port mapping does http-gateway, coaps-gatewat, UI, NATS and resource aggregate
docker run -e FQDN=172.20.1.211 -v $HOME/plgd/data:/data -d --rm --name mf_cloud -p 443:443 -p 5683:5683 -p 5684:5684 -p 9083:9083 -p 10001:10001 plgd/bundle:v2next
