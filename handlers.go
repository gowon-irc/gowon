package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ergochat/irc-go/ircevent"
	"github.com/ergochat/irc-go/ircmsg"
	"github.com/gowon-irc/go-gowon"
)

func createIRCHandler(c mqtt.Client, topic string) func(event ircmsg.Message) {
	return func(event ircmsg.Message) {
		nuh, err := ircmsg.ParseNUH(event.Source)
		if err != nil {
			log.Println(err)
		}

		line, err := event.Line()
		if err != nil {
			log.Println(err)
		}

		var msg, dest, command, args string

		if event.Command == "PRIVMSG" {
			msg = event.Params[1]
			dest = event.Params[0]
			command = gowon.GetCommand(event.Params[1])
			args = gowon.GetArgs(event.Params[1])
		}

		m := &gowon.Message{
			Module:    "gowon",
			Nick:      event.Nick(),
			Code:      event.Command,
			Raw:       line,
			Host:      nuh.Host,
			Source:    event.Source,
			User:      nuh.User,
			Arguments: event.Params,
			Tags:      event.AllTags(),
			Msg:       msg,
			Dest:      dest,
			Command:   command,
			Args:      args,
		}
		mj, err := json.Marshal(m)
		if err != nil {
			log.Println(err)

			return
		}

		c.Publish(topic, 0, false, mj)
	}
}

func createMessageHandler(irccon *ircevent.Connection, filters []string) mqtt.MessageHandler {
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

func createSendRawHandler(irccon *ircevent.Connection) mqtt.MessageHandler {
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

func createOnConnectHandler(irccon *ircevent.Connection, filters []string, topicRoot string) func(mqtt.Client) {
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
