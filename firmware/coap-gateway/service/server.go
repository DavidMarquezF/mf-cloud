package service

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/DavidMarquezF/mf-cloud/firmware/coap-gateway/uri"
	piondtls "github.com/pion/dtls/v2"
	"github.com/plgd-dev/go-coap/v2/dtls"
	"github.com/plgd-dev/go-coap/v2/mux"
	"github.com/plgd-dev/go-coap/v2/net"
	"github.com/plgd-dev/go-coap/v2/net/blockwise"
)

//Server a configuration of coapgateway
type Server struct {
	FQDN                            string // fully qualified domain name of GW
	ExternalPort                    uint16 // used to construct oic/res response
	Addr                            string // Address to listen on, ":COAP" if empty.
	IsTLSListener                   bool
	KeepaliveTimeoutConnection      time.Duration
	DisableTCPSignalMessageCSM      bool
	DisablePeerTCPSignalMessageCSMs bool
	SendErrorTextInResponse         bool
	RequestTimeout                  time.Duration
	BlockWiseTransfer               bool
	BlockWiseTransferSZX            blockwise.SZX
	ReconnectInterval               time.Duration
	HeartBeat                       time.Duration
	MaxMessageSize                  int
	LogMessages                     bool

	coapServer *dtls.Server
	listener   dtls.Listener

	ctx    context.Context
	cancel context.CancelFunc

	sigs chan os.Signal

	requestHandler RequestHandler
}

type RequestHandler struct {
	mongoClient *mongo.Client
}

type ListenCertManager = interface {
	GetServerDTLSConfig() *piondtls.Config
}

// New creates server.
func New(config Config, dtlsConfig *piondtls.Config, mongoClient *mongo.Client) *Server {
	var listener dtls.Listener

	l, err := net.NewDTLSListener("udp", config.Addr, dtlsConfig)
	if err != nil {
		log.Fatalf("cannot setup udp-dtls for server: %v", err)
	}

	listener = l

	var blockWiseTransferSZX blockwise.SZX
	switch strings.ToLower(config.BlockWiseTransferSZX) {
	case "16":
		blockWiseTransferSZX = blockwise.SZX16
	case "32":
		blockWiseTransferSZX = blockwise.SZX32
	case "64":
		blockWiseTransferSZX = blockwise.SZX64
	case "128":
		blockWiseTransferSZX = blockwise.SZX128
	case "256":
		blockWiseTransferSZX = blockwise.SZX256
	case "512":
		blockWiseTransferSZX = blockwise.SZX512
	case "1024":
		blockWiseTransferSZX = blockwise.SZX1024
	case "bert":
		blockWiseTransferSZX = blockwise.SZXBERT
	default:
		log.Fatalf("invalid value BlockWiseTransferSZX %v", config.BlockWiseTransferSZX)
	}

	ctx, cancel := context.WithCancel(context.Background())

	s := Server{
		FQDN:                       config.FQDN,
		ExternalPort:               config.ExternalPort,
		Addr:                       config.Addr,
		RequestTimeout:             config.RequestTimeout,
		SendErrorTextInResponse:    config.SendErrorTextInResponse,
		BlockWiseTransfer:          !config.DisableBlockWiseTransfer,
		BlockWiseTransferSZX:       blockWiseTransferSZX,
		ReconnectInterval:          config.ReconnectInterval,
		HeartBeat:                  config.HeartBeat,
		MaxMessageSize:             config.MaxMessageSize,
		LogMessages:                config.LogMessages,
		KeepaliveTimeoutConnection: config.KeepaliveTimeoutConnection,

		listener: listener,

		sigs: make(chan os.Signal, 1),

		ctx:    ctx,
		cancel: cancel,

		requestHandler: RequestHandler{mongoClient: mongoClient},
	}

	s.setupCoapServer()

	return &s
}

func (server *Server) setupCoapServer() {

	m := mux.NewRouter()
	//m.DefaultHandle(mux.HandlerFunc(handleA))
	m.HandleFunc(uri.Executable, requestHandler.getExecFile)

	opts := make([]dtls.ServerOption, 0, 5)
	opts = append(opts, dtls.WithMux(m))
	server.coapServer = dtls.NewServer(opts...)

}

func (server *Server) Serve() error {
	return server.serveWithHandlingSignal()
}

func (server *Server) serveWithHandlingSignal() error {
	var wg sync.WaitGroup
	var err error
	wg.Add(1)
	go func(server *Server) {
		defer wg.Done()
		err = server.coapServer.Serve(server.listener)
		server.cancel()
		server.listener.Close()
	}(server)

	signal.Notify(server.sigs,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	<-server.sigs

	server.coapServer.Stop()
	wg.Wait()

	return err
}

// Shutdown turn off server.
func (server *Server) Shutdown() error {
	select {
	case server.sigs <- syscall.SIGTERM:
	default:
	}
	return nil
}
