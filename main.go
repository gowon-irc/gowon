package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/flowchartsman/retry"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jessevdk/go-flags"

	"github.com/ergochat/irc-go/ircevent"
	"github.com/ergochat/irc-go/ircmsg"
)

var validate *validator.Validate

func setupRouter(cm *ConfigManager, cr *CommandRouter, configDir string) error {
	if err := cm.LoadDirectory(configDir); err != nil {
		return err
	}

	err := cm.Merge()
	if err != nil {
		return err
	}

	cfg := cm.MergedConfig

	validate = validator.New(validator.WithRequiredStructEnabled())

	if err := validate.RegisterValidation("irc_channel", validateIrcChannel); err != nil {
		return err
	}

	if err := validate.Struct(cfg); err != nil {
		return err
	}

	cr.Clear()

	for _, c := range cfg.Commands {
		cr.Add(&c)
	}
	cr.AddInternal("h", "list and describe commands", createHelpCommandFunc(cr))
	cr.AddInternal("gowon", "list and describe commands", createHelpCommandFunc(cr))
	cr.SortPriority()

	return nil
}

func main() {
	log.Println("starting gowon")

	opts := Config{}

	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	cm := NewConfigManager()
	cr := &CommandRouter{}

	cm.AddOpts(opts)

	if err := setupRouter(cm, cr, opts.ConfigDir); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cfg := cm.MergedConfig

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// if !event.Has(fsnotify.Chmod) {
				log.Printf("Config file %s has changed, reloading command router", event.Name)

				if err := setupRouter(cm, cr, opts.ConfigDir); err != nil {
					log.Println(err)
				}
				// }
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add(opts.ConfigDir)
	if err != nil {
		log.Fatal(err)
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
