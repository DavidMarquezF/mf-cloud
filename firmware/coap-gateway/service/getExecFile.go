package service

import (
	"bytes"
	"context"
	"io/ioutil"
	"log"

	"github.com/DavidMarquezF/mf-cloud/firmware/coap-gateway/uri"
	mfModels "github.com/DavidMarquezF/mf-cloud/firmware/mongodb"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func getData(client *mongo.Client, deviceId string) error {
	db := client.mongoClient.Database("fmw")
	coll := db.Collection("info")
	execColl := db.Collection("db")

	id, _ := primitive.ObjectIDFromHex(deviceId)

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
	log.Print(fmwInfo)

}

func (h *RequestHandler) getExecFile(w mux.ResponseWriter, r *mux.Message) {
	firmwareID := r.RouteParams.Vars[uri.FirmwareIdKey]
	log.Println(firmwareID)

	content, err := ioutil.ReadFile("test")
	err = w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader(content))
	if err != nil {
		log.Println("cannot set response: %v", err)
	}

}
