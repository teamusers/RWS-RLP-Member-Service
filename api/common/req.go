package common

import "math/big"

const (
	TYPE_CHAT_INITIAL = "chat_init"
	TYPE_CHAT_APPEND  = "chat_follow"
	METHOD_GPT        = "chatGPT"

	CODE_DIRECTION_IN  = "1"
	CODE_DIRECTION_OUT = "2"
)

type Request struct {
	Type      string `json:"type"`
	Method    string `json:"method"`
	Timestamp int64  `json:"timestamp"`
	Ascode    string `json:"ascode"`
	UserId    uint64 `json:"userid"`
	Data      string `json:"data"`
}

type OpConfig struct {
	UnitPrice *big.Int `json:"unit_price"`
	UnitLimit *big.Int `json:"unit_limit"`
	Rpc       string   `json:"rpc"`
}
