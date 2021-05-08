module github.com/DavidMarquezF/mf-cloud/firmware

go 1.14

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/miekg/dns v1.1.41 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pion/dtls/v2 v2.0.9
	github.com/plgd-dev/go-coap/v2 v2.4.0
	github.com/plgd-dev/kit v0.0.0-20210322121129-fa0d31a13679
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect

)

replace github.com/plgd-dev/go-coap/v2 => ../../go-coap

replace github.com/plgd-dev/kit => ../../kit
