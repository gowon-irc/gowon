package main

import (
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/flowchartsman/retry"
	"github.com/gowon-irc/go-gowon"
	"github.com/jessevdk/go-flags"
	irc "github.com/thoj/go-ircevent"
)

type Options struct {
	Server   string   `short:"s" long:"server" env:"GOWON_SERVER" required:"true" description:"IRC server:port"`
	User     string   `short:"u" long:"user" env:"GOWON_USER" required:"true" description:"Bot user"`
	Nick     string   `short:"n" long:"nick" env:"GOWON_NICK" required:"true" description:"Bot nick"`
	Password string   `short:"p" long:"password" env:"GOWON_PASSWORD" description:"Bot password"`
	Channels []string `short:"c" long:"channels" env:"GOWON_CHANNELS" env-delim:"," required:"true" description:"Channels to join"`
	UseTLS   bool     `short:"T" long:"tls" env:"GOWON_TLS" description:"Connect to irc server using tls"`
	Verbose  bool     `short:"v" long:"verbose" env:"GOWON_VERBOSE" description:"Verbose logging"`
	Debug    bool     `short:"d" long:"debug" env:"GOWON_DEBUG" description:"Debug logging"`
	Prefix   string   `short:"P" long:"prefix" env:"GOWON_PREFIX" default:"." description:"prefix for commands"`
	Broker   string   `short:"b" long:"broker" env:"GOWON_BROKER" default:"localhost:1883" description:"mqtt broker"`
	Filters  []string `short:"f" long:"filters" env:"GOWON_FILTERS" env-delim:"," description:"filters to apply"`
}

const (
	mqttConnectRetryInternal = 5
	mqttDisconnectTimeout    = 1000
)

func main() {
	log.Println("starting gowon")

	opts := Options{}

	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range opts.Filters {
		err := gowon.CheckFilter(f)
		if err != nil {
			log.Fatal(err)
		}
	}

	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.AddBroker(fmt.Sprintf("tcp://%s", opts.Broker))
	mqttOpts.SetClientID("gowon")
	mqttOpts.SetConnectRetry(true)
	mqttOpts.SetConnectRetryInterval(mqttConnectRetryInternal * time.Second)
	mqttOpts.SetAutoReconnect(true)

	mqttOpts.DefaultPublishHandler = defaultPublishHandler
	mqttOpts.OnConnectionLost = onConnectionLostHandler
	mqttOpts.OnReconnecting = onRecconnectingHandler

	irccon := irc.IRC(opts.Nick, opts.User)
	irccon.VerboseCallbackHandler = opts.Verbose
	irccon.Debug = opts.Debug
	irccon.UseTLS = opts.UseTLS

	if opts.Password != "" {
		irccon.Password = opts.Password
	}

	irccon.AddCallback("001", func(e *irc.Event) {
		for _, channel := range opts.Channels {
			irccon.Join(channel)
		}
	})

	mqttOpts.OnConnect = createOnConnectHandler(irccon, opts.Filters)
	c := mqtt.NewClient(mqttOpts)

	ircHandler := createIRCHandler(c)
	irccon.AddCallback("PRIVMSG", ircHandler)

	retrier := retry.NewRetrier(5, 100*time.Millisecond, 5*time.Second)
	err = retrier.Run(func() error {
		return irccon.Connect(opts.Server)
	})
	if err != nil {
		log.Fatal(err)
	}

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	irccon.Loop()

	log.Println("exiting")
	c.Disconnect(mqttDisconnectTimeout)
	log.Println("shutdown complete")
}
