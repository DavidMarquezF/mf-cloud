package service

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/DavidMarquezF/mf-cloud/firmware/coap-gateway/uri"
	mfModels "github.com/DavidMarquezF/mf-cloud/firmware/mongodb"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func getData(w mux.ResponseWriter, client *mongo.Client, deviceId string) error {
	db := client.Database("fmw")
	coll := db.Collection("info")
	execColl := db.Collection("db")

	id, _ := primitive.ObjectIDFromHex(deviceId)

	// Find the info
	filter := bson.D{{"_id", id}}

	singleResult := coll.FindOne(context.TODO(), filter)
	err := singleResult.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.SetResponse(codes.BadRequest, message.TextPlain, bytes.NewReader([]byte(err.Error())))
		} else {
			w.SetResponse(codes.InternalServerError, message.TextPlain, bytes.NewReader([]byte(err.Error())))
		}
		return err
	}
	fmwInfo := mfModels.FirmwareInfo{}

	singleResult.Decode(&fmwInfo)

	// Find the exec
	singleResultExec := execColl.FindOne(context.TODO(), bson.D{{"_id", fmwInfo.Elf}})
	err = singleResultExec.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			w.SetResponse(codes.BadRequest, message.TextPlain, strings.NewReader(fmt.Sprintf("Exec %s not found", fmwInfo.Elf)))
		} else {
			w.SetResponse(codes.InternalServerError, message.TextPlain, strings.NewReader(fmt.Sprintf("Error getting exec %s: %w", fmwInfo.Elf, err)))
		}
		return err
	}
	execInfo := mfModels.FirmwareExec{}

	singleResultExec.Decode(&execInfo)
	w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader(execInfo.Exec.Data))

	return nil
}

func (h *RequestHandler) getExecFile(w mux.ResponseWriter, r *mux.Message) {
	hashId := r.RouteParams.Vars[uri.FirmwareIdKey]
	log.Printf("Device Id: %s", hashId)
	err := getData(w, h.mongoClient, hashId)
	if err != nil {
		log.Printf("Error obtaining exec: %w", err)
		return
	}
	//content, err := ioutil.ReadFile("test")
	//err = w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader(content))
	//if err != nil {
	//		log.Println("cannot set response: %v", err)
	//}

}
