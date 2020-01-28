package server

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"time"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"

	"github.com/kolya59/virus/pkg/machine"
	pb "github.com/kolya59/virus/proto"
)

var (
	crt = "client.crt"
	key = "client.key"
	ca  = "ca.crt"
)

func getKey(uuid string) ([]byte, error) {
	// Connect to google cloud
	// Get certificate by uuid

	// Load the client certificates from disk
	certificate, err := tls.LoadX509KeyPair(crt, key)
	if err != nil {
		return nil, nil
	}

	// Create a certificate pool from the certificate authority
	certPool := x509.NewCertPool()
	ca, err := ioutil.ReadFile(ca)
	if err != nil {
		return nil, nil
	}

	// Append the certificates from the CA
	if ok := certPool.AppendCertsFromPEM(ca); !ok {
		return nil, nil
	}

	creds := credentials.NewTLS(&tls.Config{
		ServerName:   addr, // NOTE: this is required!
		Certificates: []tls.Certificate{certificate},
		RootCAs:      certPool,
	})

	return nil, nil
}

func Decode() {

}

func Encode() {

}

func SendData(machine machine.Machine) {
	// Convert to grpc
	converted := machine.ToGRPC()

	// Generate uuid
	converted.Uuid = uuid.NewV4().String()

	// TODO Get host and port
	// Set up a connection to the server.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", "host", "port"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return
	}
	defer conn.Close()
	c := pb.NewMachineSaverClient(conn)

	// Save cars
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.SaveMachine(ctx, &pb.SaveRequest{Machine: converted})
	if err != nil {
		return
	}

	// TODO: Handle error
	if r.GetMessage() == "OK" {

	}
}
