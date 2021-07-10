module github.com/DavidMarquezF/mf-cloud/firmware/coap-gateway

go 1.14

require (
	github.com/cenkalti/backoff v2.2.1+incompatible // indirect
	github.com/go-chi/chi v4.1.2+incompatible // indirect
	github.com/gofrs/uuid v3.3.0+incompatible
	github.com/gorilla/mux v1.8.0
	github.com/kelseyhightower/envconfig v1.4.0
	github.com/miekg/dns v1.1.41 // indirect
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	github.com/panjf2000/ants/v2 v2.4.3
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/pion/dtls/v2 v2.0.9
	github.com/plgd-dev/cloud v1.1.3-0.20210507093653-dea063b74876
	github.com/plgd-dev/cqrs v0.0.0-20201204150755-6ed1490c857f // indirect
	github.com/plgd-dev/go-coap/v2 v2.4.1-0.20210517130748-95c37ac8e1fa
	github.com/plgd-dev/kit v0.0.0-20210614190235-99984a49de48

	//	github.com/plgd-dev/kit v0.0.0-20210322121129-fa0d31a13679
	github.com/ugorji/go v1.1.10 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/genproto v0.0.0-20201214200347-8c77b98c765d // indirect
	google.golang.org/grpc v1.34.0
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/square/go-jose.v2 v2.5.1 // indirect
	gopkg.in/yaml.v2 v2.4.0

)

replace github.com/plgd-dev/go-coap/v2 => github.com/davidmarquezf/go-coap/v2 v2.4.1

//replace github.com/plgd-dev/kit => ../../kit

replace gopkg.in/yaml.v2 v2.3.0 => github.com/cizmazia/yaml v0.0.0-20200220134304-2008791f5454
