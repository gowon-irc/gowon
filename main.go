package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/jessevdk/go-flags"
	irc "github.com/thoj/go-ircevent"
)

type Options struct {
	Server   string   `short:"s" long:"server" env:"GOWON_SERVER" required:"true" description:"IRC server:port"`
	Nick     string   `short:"n" long:"nick" env:"GOWON_NICK" required:"true" description:"Bot nick"`
	User     string   `short:"u" long:"user" env:"GOWON_USER" required:"true" description:"Bot user"`
	Channels []string `short:"c" long:"channels" env:"GOWON_CHANNELS" required:"true" description:"Channels to join"`
	UseTLS   bool     `short:"T" long:"tls" env:"GOWON_TLS" description:"Connect to irc server using tls"`
	Verbose  bool     `short:"v" long:"verbose" env:"GOWON_VERBOSE" description:"Verbose logging"`
	Debug    bool     `short:"d" long:"debug" env:"GOWON_DEBUG" description:"Debug logging"`
	HTTPPort string   `short:"p" long:"http-port" env:"GOWON_HTTP_PORT" default:"1337" description:"http listen port"`
	Prefix   string   `short:"P" long:"prefix" env:"GOWON_PREFIX" default:"." description:"prefix for commands"`
}

func splitOptArray(sa []string) []string {
	out := []string{}

	for _, s := range sa {
		out = append(out, strings.Split(s, ",")...)
	}

	return out
}

func supCommand() string {
	return "yo"
}

func yoCommand() string {
	return "sup"
}

type httpHandler struct {
	conn *irc.Connection
}

func (h httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.conn.Privmsg("#gowon", "hello from http")

	_, err := w.Write([]byte("hello"))
	if err != nil {
		log.Println(err)
	}
}

func main() {
	opts := Options{}

	_, err := flags.Parse(&opts)
	if err != nil {
		log.Fatal(err)
	}

	channels := splitOptArray(opts.Channels)

	irccon := irc.IRC(opts.Nick, opts.User)

	irccon.VerboseCallbackHandler = opts.Verbose
	irccon.Debug = opts.Debug
	irccon.UseTLS = opts.UseTLS

	irccon.AddCallback("001", func(e *irc.Event) {
		for _, channel := range channels {
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
			}
		}(event)
	})

	h := httpHandler{irccon}

	l, err := net.Listen("tcp", fmt.Sprintf(":%s", opts.HTTPPort))
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Fatal(http.Serve(l, h))
	}()

	err = irccon.Connect(opts.Server)
	if err != nil {
		log.Fatal(err)
	}

	irccon.Loop()
}
