package main

// Simple chat using our public demo on Heroku.
// See and communicate over web version at https://jsfiddle.net/FZambia/yG7Uw/

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	centrifuge "github.com/nzlov/centrifuge-mobile"
	"github.com/nzlov/centrifugo/libcentrifugo/auth"
)

type ChatMessage struct {
	Type    string `json:"type"`
	From    string `json:"from"`
	Img     string `json:"img"`
	Name    string `json:"name"`
	Ty      string `json:"ty"`
	Content string `json:"content"`
}

// In production you need to receive credentials from application backend.
func credentials() *centrifuge.Credentials {
	// Never show secret to client of your application. Keep it on your application backend only.
	secret := "109AF84FWF45AS4S5W8F"
	// Application user ID - anonymous in this case.
	user := os.Args[2]
	// Current timestamp as string.
	timestamp := centrifuge.Timestamp()
	// Empty info.
	info := ""
	// Generate client token so Centrifugo server can trust connection parameters received from client.
	token := auth.GenerateClientToken(secret, user, timestamp, info)

	return &centrifuge.Credentials{
		User:      user,
		Timestamp: timestamp,
		Info:      info,
		Token:     token,
	}
}

type eventHandler struct {
	out io.Writer
}

func (h *eventHandler) OnConnect(c *centrifuge.Client, ctx *centrifuge.ConnectContext) {
	fmt.Fprintln(h.out, "Connected to chat")
	return
}

func (h *eventHandler) OnDisconnect(c *centrifuge.Client, ctx *centrifuge.DisconnectContext) {
	fmt.Fprintln(h.out, fmt.Sprintf("Disconnected from chat: %s", ctx.Reason))
	return
}
func (h *eventHandler) OnError(c *centrifuge.Client, ctx *centrifuge.ErrorContext) {
	fmt.Fprintln(h.out, fmt.Sprintf("Error from chat: %s", ctx))
	return
}

func (h *eventHandler) OnMessage(sub *centrifuge.Sub, msg *centrifuge.Message) {
	fmt.Fprintf(h.out, "NewMessage:%+v\n", msg)
	var chatMessage *ChatMessage
	err := json.Unmarshal(msg.Data, &chatMessage)
	if err != nil {
		return
	}
	rePrefix := fmt.Sprintf("[%v]%s says:", time.Unix(msg.Timestamp, 0), chatMessage.Name)
	fmt.Fprintln(h.out, rePrefix, chatMessage.Content)
	//sub.ReadMessage(msg.UID)
}

func (h *eventHandler) OnJoin(sub *centrifuge.Sub, info *centrifuge.ClientInfo) {
	fmt.Fprintln(h.out, fmt.Sprintf("Someone joined: user id %s", info.User))
}

func (h *eventHandler) OnRead(sub *centrifuge.Sub, msgid string) {
	fmt.Fprintln(h.out, "OnRead:", sub.Channel(), msgid)
}

func (h *eventHandler) OnLeave(sub *centrifuge.Sub, info *centrifuge.ClientInfo) {
	fmt.Fprintln(h.out, fmt.Sprintf("Someone left: user id %s", info.User))
}

func (h *eventHandler) OnSubscribeSuccess(sub *centrifuge.Sub, ctx *centrifuge.SubscribeSuccessContext) {
	fmt.Fprintln(h.out, fmt.Sprintf("Subscribed on channel %s", sub.Channel()))
}

func (h *eventHandler) OnUnsubscribe(sub *centrifuge.Sub, ctx *centrifuge.UnsubscribeContext) {
	fmt.Fprintln(h.out, fmt.Sprintf("Unsubscribed from channel %s", sub.Channel()))
}

func main() {
	creds := credentials()
	wsURL := "ws://192.168.1.9:8000/connection/websocket"

	handler := &eventHandler{os.Stdout}

	events := centrifuge.NewEventHandler()
	events.OnConnect(handler)
	events.OnError(handler)
	events.OnDisconnect(handler)
	c := centrifuge.New(wsURL, creds, events, centrifuge.DefaultConfig())

	subEvents := centrifuge.NewSubEventHandler()
	subEvents.OnMessage(handler)
	subEvents.OnRead(handler)
	subEvents.OnSubscribeSuccess(handler)
	subEvents.OnUnsubscribe(handler)
	subEvents.OnJoin(handler)
	subEvents.OnLeave(handler)

	fmt.Fprintf(os.Stdout, "You can communicate with web version at https://jsfiddle.net/FZambia/yG7Uw/\n")
	fmt.Fprintf(os.Stdout, "Connect to %s\n", wsURL)
	fmt.Fprintf(os.Stdout, "Print something and press ENTER to send\n")

	channel := os.Args[1] + ":" + os.Args[2]
	var sub *centrifuge.Sub
	var err error
	if len(os.Args) > 3 {

		sub, err = c.SubscribeWithLastMsgID(channel, os.Args[3], subEvents)
	} else {

		sub, err = c.Subscribe(channel, subEvents)
	}
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Fprintf(os.Stdout, "Print something and press ENTER to send\n")
	err = c.Connect()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Fprintf(os.Stdout, "Print something and press ENTER to send\n")

	// Read input from stdin.
	go func(sub *centrifuge.Sub) {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, _ := reader.ReadString('\n')
			msg := &ChatMessage{
				Content: text,
				Name:    "goexample",
				From:    channel,
			}
			data, _ := json.Marshal(msg)
			sub.Publish(data)
		}
	}(sub)

	// Run until CTRL+C.
	select {}
}
