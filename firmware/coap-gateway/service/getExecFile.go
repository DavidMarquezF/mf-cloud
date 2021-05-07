package service

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
)

func getExecFile(w mux.ResponseWriter, r *mux.Message) {

	content, err := ioutil.ReadFile("example.bin")
	err = w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader(content))
	if err != nil {
		log.Println("cannot set response: %v", err)
	}

}
