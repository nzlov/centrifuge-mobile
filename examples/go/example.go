package main

// Connect, subscribe on channel, publish into channel, read presence and history info.

import (
	"fmt"
	"log"

	centrifuge "github.com/nzlov/centrifuge-mobile"
	"github.com/nzlov/centrifugo/libcentrifugo/auth"
)

type TestMessage struct {
	Input string `json:"input"`
}

type subEventHandler struct{}

func (h *subEventHandler) OnMessage(sub *centrifuge.Sub, msg *centrifuge.Message) {
	log.Println(fmt.Sprintf("New message received in channel %s: %#v", sub.Channel(), msg))
	sub.ReadMessage(msg.UID)
}

func (h *subEventHandler) OnRead(sub *centrifuge.Sub, msgid string) {
	log.Println(fmt.Sprintf("New Read Message received in channel %s: %s", sub.Channel(), msgid))
}

func (h *subEventHandler) OnJoin(sub *centrifuge.Sub, msg *centrifuge.ClientInfo) {
	log.Println(fmt.Sprintf("User %s (client ID %s) joined channel %s", msg.User, msg.Client, sub.Channel()))
}

func (h *subEventHandler) OnLeave(sub *centrifuge.Sub, msg *centrifuge.ClientInfo) {
	log.Println(fmt.Sprintf("User %s (client ID %s) left channel %s", msg.User, msg.Client, sub.Channel()))
}

// In production you need to receive credentials from application backend.
func credentials() *centrifuge.Credentials {
	// Never show secret to client of your application. Keep it on your application backend only.
	secret := "109AF84FWF45AS4S5W8F"
	// Application user ID.
	user := "42"
	appkey := "web_merchant"
	// Current timestamp as string.
	timestamp := "1488055494"
	// Empty info.
	info := ""
	// Generate client token so Centrifugo server can trust connection parameters received from client.
	token := auth.GenerateClientToken(secret, user, "c31df7e10a", timestamp, info)

	fmt.Println("token:", token)
	return &centrifuge.Credentials{
		User:      user,
		Appkey:    appkey,
		Timestamp: timestamp,
		Info:      info,
		Token:     token,
	}
}

func main() {
	// In production you need to receive credentials from application backend.
	creds := credentials()

	//started := time.Now()

	wsURL := "ws://192.168.1.200:8000/connection/websocket"
	c := centrifuge.New(wsURL, creds, nil, centrifuge.DefaultConfig())
	defer c.Close()

	err := c.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	//	events := centrifuge.NewSubEventHandler()
	//	subEventHandler := &subEventHandler{}
	//	events.OnMessage(subEventHandler)
	//	events.OnJoin(subEventHandler)
	//	events.OnRead(subEventHandler)
	//	events.OnLeave(subEventHandler)
	//
	//	sub, err := c.Subscribe("public:chat", events)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//
	//	data := TestMessage{Input: fmt.Sprintf("Input:%v", time.Now())}
	//	dataBytes, _ := json.Marshal(data)
	//	_, err = sub.Publish(dataBytes)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//
	//	history, total, err := sub.History(0, -1)
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//	log.Printf("get %d messages in channel %s history,total %d", len(history), sub.Channel(), total)
	//
	//	for i, msg := range history {
	//		log.Printf("History %d : %+v\n", i, msg)
	//	}
	//
	//	presence, err := sub.Presence()
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//	log.Printf("%d clients in channel %s", len(presence), sub.Channel())
	//
	//	err = sub.Unsubscribe()
	//	if err != nil {
	//		log.Fatalln(err)
	//	}
	//
	//	log.Printf("%s", time.Since(started))
	//
	log.Println("Test Micro")

	resp, err := c.Micro("Activity.Call", `{"a":1}`)
	log.Println(string(resp.Data), err)
}
