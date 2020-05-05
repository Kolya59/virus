package models

type WSCommand struct {
	Addr  string `json:"addr"`
	Type  string `json:"type"`
	Data  []byte `json:"data"`
	Count int    `json:"count"`
}

type WSAck struct {
	Err error `json:"err,omitempty"`
}
