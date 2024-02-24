package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/ergochat/irc-go/ircevent"
	"github.com/ergochat/irc-go/ircmsg"
	"github.com/gin-gonic/gin"
	"github.com/gowon-irc/go-gowon"
)

func createIrcHandler(irccon *ircevent.Connection, cr *CommandRouter) func(event ircmsg.Message) {
	return func(event ircmsg.Message) {
		nuh, err := ircmsg.ParseNUH(event.Source)
		if err != nil {
			log.Println(err)
		}

		line, err := event.Line()
		if err != nil {
			log.Println(err)
		}

		var msg, dest, command, args string

		if event.Command == "PRIVMSG" {
			msg = event.Params[1]
			dest = event.Params[0]
			command = gowon.GetCommand(event.Params[1])
			args = gowon.GetArgs(event.Params[1])
		}

		m := &gowon.Message{
			Module:    "gowon",
			Nick:      event.Nick(),
			Code:      event.Command,
			Raw:       line,
			Host:      nuh.Host,
			Source:    event.Source,
			User:      nuh.User,
			Arguments: event.Params,
			Tags:      event.AllTags(),
			Msg:       msg,
			Dest:      dest,
			Command:   command,
			Args:      args,
		}

		rc, err := cr.Route(msg)
		if err != nil {
			return
		}

		output := rc.Send(m)

		if output == nil {
			return
		}

		for _, line := range strings.Split(output.Msg, "\n") {
			coloured := colourMsg(line)
			for _, sm := range splitMsg(coloured, 400) {
				err := irccon.Privmsg(output.Dest, sm)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}
}

func createHttpHandler(irccon *ircevent.Connection) func(*gin.Context) {
	return func(c *gin.Context) {
		var m gowon.Message

		if err := c.BindJSON(&m); err != nil {
			return
		}

		for _, line := range strings.Split(m.Msg, "\n") {
			coloured := colourMsg(line)
			for _, sm := range splitMsg(coloured, 400) {
				err := irccon.Privmsg(m.Dest, sm)
				if err != nil {
					log.Println(err)
				}
			}
		}

		c.IndentedJSON(http.StatusCreated, m)
	}
}
