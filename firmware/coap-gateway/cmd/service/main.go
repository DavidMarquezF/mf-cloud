package main

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"github.com/DavidMarquezF/mf-cloud/firmware/coap-gateway/service"
	"github.com/kelseyhightower/envconfig"
	piondtls "github.com/pion/dtls/v2"
	"github.com/plgd-dev/go-coap/v2/examples/dtls/pki"
	"github.com/plgd-dev/kit/log"
	"github.com/plgd-dev/kit/security/certManager"
	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type Config struct {
	Service service.Config
	Listen  certManager.OcfConfig `envconfig:"LISTEN"`
	Log     log.Config            `envconfig:"LOG"`
}
type Impl struct {
	service *service.Server
}

func (c Config) String() string {
	b, _ := json.MarshalIndent(c, "", "  ")
	return fmt.Sprintf("config: \n%v\n", string(b))
}

func connectToMongoDb() (*mongo.Client, error) {
	// Prepare mongodb
	mongoURI := "mongodb://" + os.Getenv("MF_MONGO_DB_SERVER") + ":27017"
	log.Printf("MONGO URI: %s", mongoURI)
	credential := options.Credential{
		Username: os.Getenv("MF_CONFIG_MONGODB_ADMINUSERNAME"),
		Password: os.Getenv("MF_CONFIG_MONGODB_ADMINPASSWORD"),
	}
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI).SetAuth(credential))
	if err != nil {
		return nil, err
	}
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}
	defer client.Disconnect(ctx)
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return nil, err
	}
	return client, nil
}
func Init(config Config) (*Impl, error) {
	log.Setup(config.Log)
	log.Info(config.String())

	dtlsConfig, err := createServerConfig(config.Listen)
	if err != nil {
		log.Fatalf("Error creating DTLS config: %v", err)
		return nil, err
	}

	client, err := connectToMongoDb()
	if err != nil {
		log.Fatalf("Error creating mongoDB connection: %v", err)
		return nil, err
	}

	return &Impl{
		service: service.New(config.Service, dtlsConfig, client),
	}, nil
}

func createServerConfig(config certManager.OcfConfig) (*piondtls.Config, error) {

	var keyBytes []byte
	var certBytes []byte
	var err error
	keyBytes, err = ioutil.ReadFile(config.File.DirPath + "/" + config.File.TLSKeyFileName)
	if err != nil {
		return nil, err
	}

	certBytes, err = ioutil.ReadFile(config.File.DirPath + "/" + config.File.TLSCertFileName)
	if err != nil {
		return nil, err
	}

	certificate, err := pki.LoadKeyAndCertificate(keyBytes, certBytes)
	if err != nil {
		return nil, err
	}
	return &piondtls.Config{
		Certificates: []tls.Certificate{*certificate},
		//	ExtendedMasterSecret: piondtls.RequireExtendedMasterSecret,
		//ClientCAs:            certPool,
		//ClientAuth: piondtls.RequireAndVerifyClientCert,
		//ConnectContextMaker: func() (context.Context, func()) {
		//		return context.WithTimeout(ctx, 30*time.Second)
		//		},
	}, nil
}

// Serve starts handling coap requests.
func (r *Impl) Serve() error {
	return r.service.Serve()
}

// Shutdown shutdowns the service.
func (r *Impl) Shutdown() error {
	err := r.service.Shutdown()
	return err
}

func main() {
	var config Config
	if err := envconfig.Process("", &config); err != nil {
		log.Fatalf("cannot parse configuration: %v", err)
	}
	if server, err := Init(config); err != nil {
		log.Fatalf("cannot init server: %v", err)
	} else {
		if err = server.Serve(); err != nil {
			log.Fatalf("unexpected ends: %v", err)
		}
	}
}
