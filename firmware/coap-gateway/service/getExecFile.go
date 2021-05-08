package service

import (
	"bytes"
	"io/ioutil"
	"log"

	"github.com/DavidMarquezF/mf-cloud/firmware/coap-gateway/uri"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
)

func getExecFile(w mux.ResponseWriter, r *mux.Message) {
	firmwareID := r.RouteParams.Vars[uri.FirmwareIdKey]
	log.Println(firmwareID)
	content, err := ioutil.ReadFile("example.bin")
	err = w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader(content))
	if err != nil {
		log.Println("cannot set response: %v", err)
	}

}
