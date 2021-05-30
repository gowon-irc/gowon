package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gowon-irc/gowon/pkg/message"
	"github.com/jessevdk/go-flags"
)

type Options struct {
	Prefix string `short:"P" long:"prefix" env:"GOWON_PREFIX" default:"." description:"prefix for commands"`
	Broker string `short:"b" long:"broker" env:"GOWON_BROKER" default:"localhost:1883" description:"mqtt broker"`
}

const mqttConnectRetryInternal = 5 * time.Second

func capHandler(m message.Message) string {
	return fmt.Sprintf("{green}%s{clear}", m.Msg)
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

	c.Subscribe("/gowon/output", 0, func(client mqtt.Client, msg mqtt.Message) {
		ms, err := message.CreateMessageStruct(msg.Payload())
		if err != nil {
			log.Print(err)

			return
		}

		if ms.Module == "module3" {
			return
		}

		var out string

		switch ms.Command {
		case "cap":
			out = capHandler(ms)
		default:
			return
		}

		ms.Module = "module3"
		ms.Msg = out
		mb, err := json.Marshal(ms)
		if err != nil {
			log.Print(err)

			return
		}
		client.Publish("/gowon/output", 0, false, mb)
	})

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
}
