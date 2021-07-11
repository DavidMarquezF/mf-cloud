module github.com/DavidMarquezF/mf-cloud/firmware/coap-gateway

go 1.14

require (
	github.com/DavidMarquezF/mf-cloud v0.0.0-20210711152338-29992d00861c
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/miekg/dns v1.1.43 // indirect
	github.com/pion/dtls/v2 v2.0.9
	github.com/plgd-dev/go-coap/v2 v2.4.1-0.20210517130748-95c37ac8e1fa
	github.com/plgd-dev/kit v0.0.0-20210614190235-99984a49de48
	go.mongodb.org/mongo-driver v1.4.2
	gopkg.in/square/go-jose.v2 v2.6.0 // indirect

)

replace github.com/plgd-dev/go-coap/v2 => github.com/davidmarquezf/go-coap/v2 v2.4.1

//replace github.com/plgd-dev/kit => ../../kit

replace gopkg.in/yaml.v2 v2.3.0 => github.com/cizmazia/yaml v0.0.0-20200220134304-2008791f5454
