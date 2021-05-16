package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/pkg/errors"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jessevdk/go-flags"
	irc "github.com/thoj/go-ircevent"
)

type Options struct {
	Server   string   `short:"s" long:"server" env:"GOWON_SERVER" required:"true" description:"IRC server:port"`
	Nick     string   `short:"n" long:"nick" env:"GOWON_NICK" required:"true" description:"Bot nick"`
	User     string   `short:"u" long:"user" env:"GOWON_USER" required:"true" description:"Bot user"`
	Channels []string `short:"c" long:"channels" env:"GOWON_CHANNELS" env-delim:"," required:"true" description:"Channels to join"`
	UseTLS   bool     `short:"T" long:"tls" env:"GOWON_TLS" description:"Connect to irc server using tls"`
	Verbose  bool     `short:"v" long:"verbose" env:"GOWON_VERBOSE" description:"Verbose logging"`
	Debug    bool     `short:"d" long:"debug" env:"GOWON_DEBUG" description:"Debug logging"`
	Prefix   string   `short:"P" long:"prefix" env:"GOWON_PREFIX" default:"." description:"prefix for commands"`
	Broker   string   `short:"b" long:"broker" env:"GOWON_BROKER" default:"localhost:1883" description:"mqtt broker"`
}

type message struct {
	Msg  string `json:"msg"`
	Nick string `json:"nick,omitempty"`
	Dest string `json:"channel"`
}

func createMessageStruct(body []byte) (m message, err error) {
	err = json.Unmarshal(body, &m)
	if err != nil {
		return m, errors.Wrap(err, "message couldn't be parsed as message json")
	}

	if m.Msg == "" {
		return m, errors.New("message body does not contain any message content")
	}

	if m.Dest == "" {
		return m, errors.New("message body does not contain a destination")
	}

	return m, nil
}

func createMessageHandler(irccon *irc.Connection) mqtt.MessageHandler {
	return func(client mqtt.Client, msg mqtt.Message) {
		m, err := createMessageStruct(msg.Payload())
		if err != nil {
			log.Print(err)
			return
		}

		irccon.Privmsg(m.Dest, m.Msg)
	}
}

func main() {
	opts := Options{}

	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	mqttOpts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s", opts.Broker))
	mqttOpts.SetClientID("gowon")

	c := mqtt.NewClient(mqttOpts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	irccon := irc.IRC(opts.Nick, opts.User)
	irccon.VerboseCallbackHandler = opts.Verbose
	irccon.Debug = opts.Debug
	irccon.UseTLS = opts.UseTLS

	irccon.AddCallback("001", func(e *irc.Event) {
		for _, channel := range opts.Channels {
			irccon.Join(channel)
		}
	})

	irccon.AddCallback("PRIVMSG", func(event *irc.Event) {
		go func(event *irc.Event) {
			m := &message{
				Dest: event.Arguments[0],
				Msg:  event.Arguments[1],
				Nick: event.Nick,
			}

			msgJSON, _ := json.Marshal(m)

			c.Publish("/gowon/input", 0, false, msgJSON)
		}(event)
	})

	msgHandler := createMessageHandler(irccon)
	c.Subscribe("/gowon/output", 0, msgHandler)

	err = irccon.Connect(opts.Server)
	if err != nil {
		log.Fatal(err)
	}

	irccon.Loop()
}
