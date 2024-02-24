package main

import (
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gowon-irc/go-gowon"
)

const (
	moduleName = "module1"
)

func main() {
	r := gin.Default()
	r.POST("/message", func(c *gin.Context) {
		var msg gowon.Message

		if err := c.BindJSON(&msg); err != nil {
			log.Printf("Error: unable to bind message to json: %v", err)
			return
		}

		returnMsg := &gowon.Message{
			Module: moduleName,
			Msg:    strings.ToUpper(msg.Args),
			Dest:   msg.Dest,
		}

		c.IndentedJSON(http.StatusOK, returnMsg)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
