package service

import (
	"bytes"
	"io/ioutil"
	"log"
)

func getExecFile(w mux.ResponseWriter, r *mux.Message) {

	content, err := ioutil.ReadFile("test.txt")
	err = w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader(content))
	if err != nil {
		log.Println("cannot set response: %v", err)
	}

}
