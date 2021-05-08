package main

import (
	"context"
	"log"
	"os"
	"time"

	piondtls "github.com/pion/dtls/v2"
	"github.com/plgd-dev/go-coap/v2/dtls"
	"github.com/plgd-dev/go-coap/v2/net/blockwise"
)

func main() {
	// root cert

	config := &piondtls.Config{
		//Certificates: []tls.Certificate{},
		//ExtendedMasterSecret: piondtls.NoClientCert,
		//RootCAs:            nil,
		InsecureSkipVerify: true,
	}
	co, err := dtls.Dial("localhost:5682", config, dtls.WithBlockwise(true, blockwise.SZX16, time.Second*10))
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	path := "/api/v1/fmw/asd/exec"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	resp, err := co.Get(ctx, path)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	log.Printf("Response payload: %v", resp.String())
	m, _ := resp.Marshal()
	log.Print(string(m))
}
