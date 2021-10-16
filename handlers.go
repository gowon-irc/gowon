package main

import (
	"encoding/json"
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gowon-irc/go-gowon"
	irc "github.com/thoj/go-ircevent"
)

func createIRCHandler(c mqtt.Client) func(event *irc.Event) {
	return func(event *irc.Event) {
		go func(event *irc.Event) {
			m := &gowon.Message{
				Module: "gowon",
				Dest:   event.Arguments[0],
				Msg:    event.Arguments[1],
				Nick:   event.Nick,
			}
			mj, err := json.Marshal(m)
			if err != nil {
				log.Print(err)

				return
			}

			c.Publish("/gowon/input", 0, false, mj)
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

func defaultPublishHandler(c mqtt.Client, msg mqtt.Message) {
	log.Printf("unexpected message:  %s\n", msg)
}

func onConnectionLostHandler(c mqtt.Client, err error) {
	log.Println("connection to broker lost")
}

func onRecconnectingHandler(c mqtt.Client, opts *mqtt.ClientOptions) {
	log.Println("attempting to reconnect to broker")
}

func createOnConnectHandler(irccon *irc.Connection, filters []string) func(mqtt.Client) {
	log.Println("connected to broker")

	mh := createMessageHandler(irccon, filters)

	return func(client mqtt.Client) {
		client.Subscribe("/gowon/output", 0, mh)

		log.Printf("Subscription to /gowon/output complete")
	}
}
