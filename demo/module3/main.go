package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gowon-irc/go-gowon"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Prefix string `short:"P" long:"prefix" env:"GOWON_PREFIX" default:"." description:"prefix for commands"`
	Broker string `short:"b" long:"broker" env:"GOWON_BROKER" default:"localhost:1883" description:"mqtt broker"`
}

const mqttConnectRetryInternal = 5 * time.Second

func capHandler(m gowon.Message) (string, error) {
	filtered, _ := gowon.Filter(&m, "module=module3")
	if filtered {
		return "", nil
	}

	return fmt.Sprintf("{green}%s{clear}", m.Msg), nil
}

func main() {
	opts := Options{}

	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	mqttOpts := mqtt.NewClientOptions().AddBroker(fmt.Sprintf("tcp://%s", opts.Broker))
	mqttOpts.SetClientID("gowon_module3")
	mqttOpts.SetConnectRetry(true)
	mqttOpts.SetConnectRetryInterval(mqttConnectRetryInternal)

	c := mqtt.NewClient(mqttOpts)
	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	mr := gowon.NewMessageRouter()
	mr.AddCommand("cap", capHandler)
	mr.SubscribeMiddleware(c, "module3")

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}
