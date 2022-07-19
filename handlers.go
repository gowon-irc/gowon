package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gowon-irc/go-gowon"
	irc "github.com/thoj/go-ircevent"
)

func createIRCHandler(c mqtt.Client, topic string) func(event *irc.Event) {
	return func(event *irc.Event) {
		go func(event *irc.Event) {
			var msg, dest, command, args string

			if event.Code == "PRIVMSG" {
				msg = event.Arguments[1]
				dest = event.Arguments[0]
				command = gowon.GetCommand(event.Arguments[1])
				args = gowon.GetArgs(event.Arguments[1])
			}

			m := &gowon.Message{
				Module:    "gowon",
				Nick:      event.Nick,
				Code:      event.Code,
				Raw:       event.Raw,
				Host:      event.Host,
				Source:    event.Host,
				User:      event.User,
				Arguments: event.Arguments,
				Tags:      event.Tags,
				Msg:       msg,
				Dest:      dest,
				Command:   command,
				Args:      args,
			}
			mj, err := json.Marshal(m)
			if err != nil {
				log.Print(err)

				return
			}

			c.Publish(topic, 0, false, mj)
		}(event)
	}
}

func createMessageHandler(irccon *irc.Connection, filters []string) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		m, err := gowon.CreateMessageStruct(msg.Payload())
		if err != nil {
			log.Print(err)

			return
		}

		for _, f := range filters {
			filtered, err := gowon.Filter(&m, f)

			if err != nil {
				break
			}

			if filtered {
				log.Printf(`Message "%s" has been filtered by filter "%s"`, m.Msg, f)
				return
			}
		}

		for _, line := range strings.Split(m.Msg, "\n") {
			coloured := colourMsg(line)
			for _, sm := range splitMsg(coloured, 400) {
				irccon.Privmsg(m.Dest, sm)
			}
		}
	}
}

func createSendRawHandler(irccon *irc.Connection) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		irccon.SendRaw(string(msg.Payload()))
	}
}

func defaultPublishHandler(c mqtt.Client, msg mqtt.Message) {
	log.Printf("unexpected message:  %s\n", msg)
}

func onConnectionLostHandler(c mqtt.Client, err error) {
	log.Println("connection to broker lost")
}

func onRecconnectingHandler(c mqtt.Client, opts *mqtt.ClientOptions) {
	log.Println("attempting to reconnect to broker")
}

func createOnConnectHandler(irccon *irc.Connection, filters []string, topicRoot string) func(mqtt.Client) {
	log.Println("connected to broker")

	topic := topicRoot + "/output"
	rawTopic := topicRoot + "/raw/output"

	mh := createMessageHandler(irccon, filters)
	rh := createSendRawHandler(irccon)

	return func(client mqtt.Client) {
		client.Subscribe(topic, 0, mh)
		log.Printf(fmt.Sprintf("Subscription to %s complete", topic))

		client.Subscribe(rawTopic, 0, rh)
		log.Printf(fmt.Sprintf("Subscription to %s complete", rawTopic))
	}
}
