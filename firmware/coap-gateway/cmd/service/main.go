package main

import (
	"fmt"

	"github.com/DavidMarquezF/mf-cloud/firmware/coap-gateway/service"
	"github.com/kelseyhightower/envconfig"
	"github.com/plgd-dev/kit/log"
	"github.com/plgd-dev/kit/security/certManager"
)

type Config struct {
	Service service.Config
	Listen  certManager.OcfConfig `envconfig:"LISTEN"`
	Log     log.Config            `envconfig:"LOG"`
}
type Impl struct {
	service           *service.Server
}

func Init(config Config) (*Impl, error) {
	log.Setup(config.Log)
	log.Info(config.String())

	dtlsConfig, err := createServerConfig(config.Listen)
	if err != nil{
		log.Fatalf("Error creating DTLS config: %v", err)
		return nil, err
	}

	return &Impl{
		service:           service.New(config.Service, dtlsConfig),
	}, nil
}

func createServerConfig(config *certManager.OcfConfig) (*piondtls.Config, error) {

	var keyBytes []byte
	var certBytes []byte
	var err error
	keyBytes, err = ioutil.ReadFile(keyPath := config.DirPath + "/" + config.TLSKeyFileName)
	if err != nil {
		return nil, err
	}

	certBytes, err = ioutil.ReadFile(config.DirPath + "/" + config.TLSCertFileName)
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
func (r *RefImpl) Serve() error {
	return r.service.Serve()
}

// Shutdown shutdowns the service.
func (r *RefImpl) Shutdown() error {
	err := r.service.Shutdown()
	r.dialCertManager.Close()
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
