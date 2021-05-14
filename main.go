package main

import (
	"fmt"
	"log"

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

func supCommand() string {
	return "yo"
}

func yoCommand() string {
	return "sup"
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
			channel := event.Arguments[0]
			command := event.Arguments[1]

			cm := map[string]func() string{
				opts.Prefix + "sup": supCommand,
				opts.Prefix + "yo":  yoCommand,
			}

			if f, ok := cm[command]; ok {
				irccon.Privmsg(channel, f())
				c.Publish("/test/a", 0, false, f())
			}
		}(event)
	})

	err = irccon.Connect(opts.Server)
	if err != nil {
		log.Fatal(err)
	}

	irccon.Loop()
}
