package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gowon-irc/gowon/pkg/message"
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
	Filters  []string `short:"f" long:"filters" env:"GOWON_FILTERS" env-delim:"," description:"filters to apply"`
}

const mqttConnectRetryInternal = 5 * time.Second

func createIRCHandler(c mqtt.Client) func(event *irc.Event) {
	return func(event *irc.Event) {
		go func(event *irc.Event) {
			m := &message.Message{
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
		m, err := message.CreateMessageStruct(msg.Payload())
		if err != nil {
			log.Print(err)

			return
		}

		for _, f := range filters {
			filtered, err := message.Filter(&m, f)

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

func main() {
	opts := Options{}

	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range opts.Filters {
		err := message.CheckFilter(f)
		if err != nil {
			log.Fatal(err)
		}
	}

	mqttOpts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s", opts.Broker))
	mqttOpts.SetClientID("gowon")
	mqttOpts.SetConnectRetry(true)
	mqttOpts.SetConnectRetryInterval(mqttConnectRetryInternal)

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

	ircHandler := createIRCHandler(c)
	irccon.AddCallback("PRIVMSG", ircHandler)

	msgHandler := createMessageHandler(irccon, opts.Filters)
	c.Subscribe("/gowon/output", 0, msgHandler)

	err = irccon.Connect(opts.Server)
	if err != nil {
		log.Fatal(err)
	}

	irccon.Loop()
}
