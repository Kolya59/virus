package models

type WSCommand struct {
	Addr  string `json:"addr"`
	Type  string `json:"type"`
	Data  []byte `json:"data"`
	Count int    `json:"count"`
}

type WSAck struct {
	Err string `json:"err,omitempty"`
}
