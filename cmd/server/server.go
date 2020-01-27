package main

import (
	"encoding/json"
	"fmt"

	"github.com/kolya59/virus/pkg/machine"
)

func main() {
	m := machine.Machine{}
	m.GetIPS()
	j, _ := json.Marshal(&m)
	fmt.Print(string(j))
}
