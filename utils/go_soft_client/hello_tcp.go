package main

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/plgd-dev/go-coap/v2/net/blockwise"

	"github.com/plgd-dev/go-coap/v2/tcp"
)

func main() {
	co, err := tcp.Dial("localhost:5682", /* tcp.WithTLS(&tls.Config{
			InsecureSkipVerify: true,
		}),*/tcp.WithBlockwise(true, blockwise.SZX1024, time.Second*10), tcp.WithMaxMessageSize(int(blockwise.SZX1024.Size())))
	if err != nil {
		log.Fatalf("Error dialing: %v", err)
	}
	path := "/api/v1/fmw/{firmwareid}/exec"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	resp, err := co.Get(ctx, path)
	if err != nil {
		log.Fatalf("Error sending request: %v", err)
	}
	log.Printf("Response payload: %v", resp.String())
	m, _ := resp.Marshal()
	log.Print(string(m))
}
