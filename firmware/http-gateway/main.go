package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"html"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	mfModules "github.com/DavidMarquezF/mf-cloud/firmware/modules"
	"github.com/plgd-dev/kit/codec/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	API string = "/api/v1"

	//firmware
	Firmware = API + "/fmw"
)

type SoftwareResource struct {
	Purl           string `json:"purl"`
	Swupdateaction string `json:"swupdateaction"`
	Updatetime     string `json:"updatetime"`
}

type CreateFirmwareRequest struct {
	DeviceId  string               `json:"device_id"`
	PlgdCloud string               `json:"plgd_cloud"`
	Token     string               `json:"token"` // Don't do this ever!! This is just to try out the idea
	Modules   []mfModules.ModuleId `json:"modules"`
}

type BaseHandler struct {
	db  *mongo.Client
	ctx *context.Context
}

func errToJsonRes(err error) map[string]string {
	return map[string]string{"err": err.Error()}
}

func writeError(w http.ResponseWriter, err error) {
	if err == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	b, _ := json.Encode(errToJsonRes(err))
	http.Error(w, string(b), http.StatusBadRequest)
}

func changePurl(sw *CreateFirmwareRequest, w http.ResponseWriter, url string) error {
	swUpdatebody := SoftwareResource{
		Purl:           url,
		Swupdateaction: "idle",
		Updatetime:     "2030-01-01T00:00:00Z",
	}
	b, _ := json.Encode(swUpdatebody)
	req, err := http.NewRequest(http.MethodPut, sw.PlgdCloud+"/api/v1/devices/"+sw.DeviceId+"/sw", bytes.NewReader(b))

	// add authorization header to the req
	req.Header.Add("Authorization", sw.Token)

	// Send req using http Client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		writeError(w, fmt.Errorf("Cannot change purl: %w", err))
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		writeError(w, fmt.Errorf("Cannot change purl: %w", err))
		return err
	}

	statusOk := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOk {
		http.Error(w, string(body), resp.StatusCode)
		return fmt.Errorf("Cannot change purl: %w", err)
	}

	log.Println(string([]byte(body)))

	return nil
}

func updateDevice(sw *CreateFirmwareRequest, w http.ResponseWriter) error {

	req, err := http.NewRequest(http.MethodPut, sw.PlgdCloud+"/api/v1/devices/"+sw.DeviceId+"/oic/sec/pstat", bytes.NewBufferString(`{"tm": 256}`))

	// add authorization header to the req
	req.Header.Add("Authorization", sw.Token)

	// Send req using http Client
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Do(req)
	if err != nil {
		writeError(w, fmt.Errorf("Cannot update device: %w", err))
		return err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		writeError(w, fmt.Errorf("Cannot change update device: %w", err))
		return err
	}

	statusOk := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOk {
		http.Error(w, string(body), resp.StatusCode)
		return errors.New("error")
	}

	log.Println(string([]byte(body)))

	return nil
}

type FirmwareInfo struct {
	ID      primitive.ObjectID `bson:"_id,omitempty"`
	Date    primitive.DateTime `bson:"date"`
	Version string             `bson:"version"`
	Elf     primitive.ObjectID `bson:"elf"`
}

type FirmwareExec struct {
	ID   primitive.ObjectID `bson:"_id,omitempty"`
	Exec primitive.Binary   `bson:"exec"`
}

func (h *BaseHandler) createFirmware(w http.ResponseWriter, r *http.Request) {

	var body CreateFirmwareRequest
	if err := json.ReadFrom(r.Body, &body); err != nil {
		writeError(w, fmt.Errorf("invalid json body: %w", err))
		return
	}

	//updateDevice(&body, w)
	modules := make([]mfModules.Module, len(body.Modules))
	for i, s := range body.Modules {
		modules[i] = mfModules.Module{Id: s}
	}

	conf := mfModules.FirmwareConfig{
		DeviceId:   body.DeviceId,
		DeviceName: "Mf Device",
		Platform:   mfModules.ESP32,
		Modules:    modules,
	}

	b, _ := json.Encode(conf)
	req, err := http.NewRequest(http.MethodPut, "http://"+os.Getenv("MF_BUILDER_SERVER")+":"+os.Getenv("MF_BUILDER_SERVER_PORT")+"/create-firmware", bytes.NewReader(b))

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		writeError(w, fmt.Errorf("error connecting to builder server: %w", err))
		return
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	statusOk := resp.StatusCode >= 200 && resp.StatusCode < 300
	if !statusOk {
		http.Error(w, `{err: "Error in build server", innerErr: `+string(respBody)+`}`, resp.StatusCode)
		return
	}

	db := h.db.Database("fmw")
	coll := db.Collection("info")
	execColl := db.Collection("db")

	execData := FirmwareExec{
		Exec: primitive.Binary{Data: respBody},
	}
	execInsertResult, err := execColl.InsertOne(context.TODO(), execData)
	if err != nil {
		writeError(w, fmt.Errorf("Cannot create db entry for elf: %w", err))
		return
	}
	execInsertedId := execInsertResult.InsertedID.(primitive.ObjectID)

	filter := bson.D{{"_id", body.DeviceId}}
	id, _ := primitive.ObjectIDFromHex(body.DeviceId)

	newData := FirmwareInfo{
		ID:      id,
		Version: "1.0.0",
		Date:    primitive.NewDateTimeFromTime(time.Now()),
		Elf:     execInsertedId,
	}
	oldData := FirmwareInfo{}
	singleResult := coll.FindOneAndReplace(context.TODO(), filter, newData)
	err = singleResult.Err()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			_, err := coll.InsertOne(context.TODO(), newData)
			if err != nil {
				writeError(w, fmt.Errorf("Cannot insert info db entry: %w", err))
				return
			}

		} else {
			writeError(w, fmt.Errorf("Cannot update db entry for info: %w", err))
			return
		}

	}

	if singleResult.Err() != mongo.ErrNoDocuments {
		singleResult.Decode(&oldData)
		_, err = execColl.DeleteOne(context.TODO(), bson.M{"_id": oldData})
		if err != nil {
			writeError(w, fmt.Errorf("Error deleting old exec: %w", err))
			return
		}
	}

	err = changePurl(body, w, "coap://"+os.Getenv("MF_COAP_GATEWAY_SERVER")+":"+os.Getenv("MF_COAP_GATEWAY_SERVER_PORT")+"/api/v1/fmw/"+body.DeviceId+"/exec")
	if err != nil {
		return
	}
	err = updateDevice()
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonData := []byte(`{"status":"OK"}`)
	w.Write(jsonData)

}

func respondAlive(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
}

func newBaseHandler(db *mongo.Client, context *context.Context) *BaseHandler {
	return &BaseHandler{db: db, ctx: context}
}

func main() {
	// Prepare mongodb
	mongoURI := "mongodb://" + os.Getenv("MF_MONGO_DB_SERVER") + ":27017"
	log.Printf("MONGO URI: %s", mongoURI)
	credential := options.Credential{
		Username: os.Getenv("MF_CONFIG_MONGODB_ADMINUSERNAME"),
		Password: os.Getenv("MF_CONFIG_MONGODB_ADMINPASSWORD"),
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI).SetAuth(credential))
	if err != nil {
		log.Fatal(err)
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}

	h := newBaseHandler(client, &ctx)
	// Prepare server
	http.HandleFunc(Firmware, h.createFirmware)
	http.HandleFunc("/", respondAlive)

	http.ListenAndServe(":8090", nil)
}
