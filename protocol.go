package centrifuge

import (
	"encoding/json"
)

type clientCommand struct {
	UID    string `json:"uid"`
	Method string `json:"method"`
}

type connectClientCommand struct {
	clientCommand
	Params connectParams `json:"params"`
}

type refreshClientCommand struct {
	clientCommand
	Params refreshParams `json:"params"`
}

type subscribeClientCommand struct {
	clientCommand
	Params subscribeParams `json:"params"`
}

type unsubscribeClientCommand struct {
	clientCommand
	Params unsubscribeParams `json:"params"`
}

type publishClientCommand struct {
	clientCommand
	Params publishParams `json:"params"`
}

type presenceClientCommand struct {
	clientCommand
	Params presenceParams `json:"params"`
}

type readClientCommand struct {
	clientCommand
	Params readParams `json:"params"`
}

type historyClientCommand struct {
	clientCommand
	Params historyParams `json:"params"`
}

type pingClientCommand struct {
	clientCommand
}

type connectParams struct {
	User      string `json:"user"`
	Timestamp string `json:"timestamp"`
	Info      string `json:"info"`
	Token     string `json:"token"`
}

type refreshParams struct {
	User      string `json:"user"`
	Timestamp string `json:"timestamp"`
	Info      string `json:"info"`
	Token     string `json:"token"`
}

type subscribeParams struct {
	Channel string `json:"channel"`
	Client  string `json:"client,omitempty"`
	Last    string `json:"last,omitempty"`
	Recover bool   `json:"recover,omitempty"`
	Info    string `json:"info,omitempty"`
	Sign    string `json:"sign,omitempty"`
}

type unsubscribeParams struct {
	Channel string `json:"channel"`
}

type publishParams struct {
	Channel string           `json:"channel"`
	Data    *json.RawMessage `json:"data"`
}

type presenceParams struct {
	Channel string `json:"channel"`
}

type historyParams struct {
	Channel string `json:"channel"`
	Skip    int    `json:"skip"`
	Limit   int    `json:"limit"`
}

type readParams struct {
	Channel string `json:"channel"`
	MsgID   string `json:"msgid"`
}

type response struct {
	UID    string          `json:"uid,omitempty"`
	Error  string          `json:"error"`
	Method string          `json:"method"`
	Body   json.RawMessage `json:"body"`
}

type joinLeaveMessage struct {
	Channel string        `json:"channel"`
	Data    rawClientInfo `json:"data"`
}

type connectResponseBody struct {
	Version string `json:"version"`
	Client  string `json:"client"`
	Expires bool   `json:"expires"`
	Expired bool   `json:"expired"`
	TTL     int64  `json:"ttl"`
}

type subscribeResponseBody struct {
	Channel   string       `json:"channel"`
	Status    bool         `json:"status"`
	Last      string       `json:"last"`
	Messages  []rawMessage `json:"messages"`
	Recovered bool         `json:"recovered"`
}

type unsubscribeResponseBody struct {
	Channel string `json:"channel"`
	Status  bool   `json:"status"`
}

type publishResponseBody struct {
	Channel string `json:"channel"`
	Status  bool   `json:"status"`
}

type presenceResponseBody struct {
	Channel string                   `json:"channel"`
	Data    map[string]rawClientInfo `json:"data"`
}

type historyResponseBody struct {
	Channel string       `json:"channel"`
	Data    []rawMessage `json:"data"`
	Total   int          `json:"total"`
}
type readResponseBody struct {
	Channel string `json:"channel"`
	MsgID   string `json:"msgid"`
	Read    bool   `json:"read"`
}

type disconnectAdvice struct {
	Reason    string `json:"reason"`
	Reconnect bool   `json:"reconnect"`
}

var (
	arrayJsonPrefix  byte = '['
	objectJsonPrefix byte = '{'
)

func responsesFromClientMsg(msg []byte) ([]response, error) {
	var resps []response
	firstByte := msg[0]
	switch firstByte {
	case objectJsonPrefix:
		// single command request
		var resp response
		err := json.Unmarshal(msg, &resp)
		if err != nil {
			return nil, err
		}
		resps = append(resps, resp)
	case arrayJsonPrefix:
		// array of commands received
		err := json.Unmarshal(msg, &resps)
		if err != nil {
			return nil, err
		}
	default:
		return nil, ErrInvalidMessage
	}
	return resps, nil
}
