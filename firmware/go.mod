module github.com/DavidMarquezF/mf-cloud/firmware

go 1.14

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/miekg/dns v1.1.41 // indirect
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pion/dtls/v2 v2.0.1-0.20200503085337-8e86b3a7d585
	github.com/plgd-dev/go-coap/v2 v2.1.4-0.20201201213140-b8c428d8fccf
	github.com/plgd-dev/kit v0.0.0-20210322121129-fa0d31a13679
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect

)

replace github.com/plgd-dev/go-coap/v2 => github.com/plgd-dev/go-coap/v2 v2.1.4-0.20201201213140-b8c428d8fccf
