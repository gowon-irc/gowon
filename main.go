package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/flowchartsman/retry"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jessevdk/go-flags"

	"github.com/ergochat/irc-go/ircevent"
	"github.com/ergochat/irc-go/ircmsg"
)

const (
	mqttConnectRetryInternal = 5
	mqttDisconnectTimeout    = 1000
	configFilename           = "config.yaml"
)

var validate *validator.Validate

func main() {
	log.Println("starting gowon")

	opts := Config{}

	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	cm := NewConfigManager()
	cm.AddOpts(opts)
	if err := cm.LoadDirectory(opts.ConfigDir); err != nil {
		log.Println(err)
	}

	cfg, err := cm.Merge()
	if err != nil {
		log.Fatal(err)
	}

	validate = validator.New(validator.WithRequiredStructEnabled())

	if err := validate.RegisterValidation("irc_channel", validateIrcChannel); err != nil {
		log.Fatalf("failed to register validation, err : %v", err)
	}

	if err := validate.Struct(cfg); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	mqttOpts := mqtt.NewClientOptions()
	mqttOpts.AddBroker(fmt.Sprintf("tcp://%s", cfg.Broker))
	mqttOpts.SetClientID("gowon")
	mqttOpts.SetConnectRetry(true)
	mqttOpts.SetConnectRetryInterval(mqttConnectRetryInternal * time.Second)
	mqttOpts.SetAutoReconnect(true)

	mqttOpts.DefaultPublishHandler = defaultPublishHandler
	mqttOpts.OnConnectionLost = onConnectionLostHandler
	mqttOpts.OnReconnecting = onRecconnectingHandler

	irccon := ircevent.Connection{
		Server:      cfg.Server,
		Nick:        cfg.Nick,
		User:        cfg.User,
		Debug:       cfg.Debug,
		RequestCaps: []string{"server-time"},
	}
	// ircevent.VerboseCallbackHandler = cfg.Verbose

	irccon.UseTLS = cfg.UseTLS
	if cfg.UseTLS {
		irccon.TLSConfig = &tls.Config{
			ServerName: strings.Split(cfg.Server, ":")[0],
			MinVersion: tls.VersionTLS12,
		}
	}

	if cfg.Password != "" {
		irccon.UseSASL = true
		irccon.SASLLogin = cfg.Nick
		irccon.SASLPassword = cfg.Password
	}

	irccon.AddConnectCallback(func(e ircmsg.Message) {
		for _, channel := range cfg.Channels {
			irccon.Join(channel)
		}
	})

	mqttOpts.OnConnect = createOnConnectHandler(&irccon, cfg.TopicRoot)
	c := mqtt.NewClient(mqttOpts)

	privMsgHandler := createIRCHandler(c, cfg.TopicRoot+"/input")
	irccon.AddCallback("PRIVMSG", privMsgHandler)

	ircRawHandler := createIRCHandler(c, cfg.TopicRoot+"/raw/input")
	// irccon.AddCallback("*", ircRawHandler)
	for _, c := range []string{"JOIN", "332", "353"} {
		irccon.AddCallback(c, ircRawHandler)
	}

	httpRouter := gin.Default()
	httpRouter.POST("/message", createHttpHandler(&irccon))

	retrier := retry.NewRetrier(5, 100*time.Millisecond, 5*time.Second)
	err = retrier.Run(func() error {
		return irccon.Connect()
	})
	if err != nil {
		log.Fatal(err)
	}

	if token := c.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	go func() {
		if err := httpRouter.Run(fmt.Sprintf("0.0.0.0:%d", cfg.HttpPort)); err != nil {
			log.Fatal(err)
		}
	}()

	irccon.Loop()

	log.Println("exiting")
	c.Disconnect(mqttDisconnectTimeout)
	log.Println("shutdown complete")
}
