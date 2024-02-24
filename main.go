package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/flowchartsman/retry"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jessevdk/go-flags"

	"github.com/ergochat/irc-go/ircevent"
	"github.com/ergochat/irc-go/ircmsg"
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
			if err = irccon.Join(channel); err != nil {
				log.Println(err)
			}
		}
	})

	c1 := &Command{Command: "c1"}
	c2 := &Command{Command: "c2"}
	cr := &CommandRouter{}
	cr.Add(c1)
	cr.Add(c2)

	privMsgHandler := createIrcHandler(&irccon, cr)
	irccon.AddCallback("PRIVMSG", privMsgHandler)

	httpRouter := gin.Default()
	httpRouter.POST("/message", createHttpHandler(&irccon))

	retrier := retry.NewRetrier(5, 100*time.Millisecond, 5*time.Second)
	err = retrier.Run(func() error {
		return irccon.Connect()
	})
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		if err := httpRouter.Run(fmt.Sprintf("0.0.0.0:%d", cfg.HttpPort)); err != nil {
			log.Fatal(err)
		}
	}()

	irccon.Loop()

	log.Println("shutdown complete")
}
