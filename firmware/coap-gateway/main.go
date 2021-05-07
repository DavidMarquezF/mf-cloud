package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"

	piondtls "github.com/pion/dtls/v2"
	"github.com/plgd-dev/go-coap/v2/dtls"
	"github.com/plgd-dev/go-coap/v2/examples/dtls/pki"
	"github.com/plgd-dev/go-coap/v2/message"
	"github.com/plgd-dev/go-coap/v2/message/codes"
	"github.com/plgd-dev/go-coap/v2/mux"
	"github.com/plgd-dev/go-coap/v2/net"
	"github.com/plgd-dev/go-coap/v2/udp/client"
)

func toHexInt(n *big.Int) string {
	return fmt.Sprintf("%x", n) // or %X or upper case
}

func handleA(w mux.ResponseWriter, r *mux.Message) {
	//	clientCert := r.Context.Value("client-cert").(*x509.Certificate)
	//	log.Println("Serial number:", toHexInt(clientCert.SerialNumber))
	//	log.Println("Subject:", clientCert.Subject)
	//	log.Println("Email:", clientCert.EmailAddresses)

	content, err := ioutil.ReadFile("test.txt")
	err = w.SetResponse(codes.Content, message.TextPlain, bytes.NewReader(content))
	if err != nil {
		log.Println("cannot set response: %v", err)
	}

}

func onNewClientConn(cc *client.ClientConn, dtlsConn *piondtls.Conn) {

}

func main() {
	m := mux.NewRouter()
	m.DefaultHandle(mux.HandlerFunc(handleA))
	m.Handle(Executable, mux.HandlerFunc(handleA))

	config, err := createServerConfig(context.Background())
	if err != nil {
		log.Fatalln(err)
		return
	}

	log.Fatal(listenAndServeDTLS("udp", ":5682", config, m))
}

func listenAndServeDTLS(network string, addr string, config *piondtls.Config, handler mux.Handler) error {
	l, err := net.NewDTLSListener(network, addr, config)
	if err != nil {
		return err
	}
	defer l.Close()
	s := dtls.NewServer(dtls.WithMux(handler), dtls.WithOnNewClientConn(onNewClientConn))
	return s.Serve(l)
}

func createServerConfig(ctx context.Context) (*piondtls.Config, error) {

	var keyBytes []byte
	var certBytes []byte
	var err error
	keyBytes, err = ioutil.ReadFile("certs/external/firmware-coap.key")
	if err != nil {
		return nil, err
	}

	certBytes, err = ioutil.ReadFile("certs/external/firmware-coap.crt")
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
