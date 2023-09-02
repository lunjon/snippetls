package types

import (
	"encoding/json"
)

type RequestMessage struct {
	JsonRpc string `json:"jsonrpc"`
	ID      int    `json:"id,omitempty"`
	Method  string `json:"method"`
	Params  any    `json:"params"`
}

func NewRequest(method string, params any) RequestMessage {
	return RequestMessage{
		JsonRpc: "2.0",
		Method:  method,
		Params:  params,
	}
}

func ParamsAs[T any](p any) (T, error) {
	var t T
	bs, err := json.Marshal(p)
	if err != nil {
		return t, err
	}

	err = json.Unmarshal(bs, &t)
	return t, err
}
