package types

type ResponseMessage struct {
	JsonRpc string `json:"jsonrpc"`
	ID      int    `json:"id"`
	Result  any    `json:"result"`
}

func NewResponse(id int, result any) ResponseMessage {
	m := ResponseMessage{
		JsonRpc: "2.0",
		ID:      id,
		Result:  result,
	}
	return m
}
