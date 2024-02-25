package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gowon-irc/go-gowon"
)

const (
	moduleName = "module4"
)

func main() {
	r := gin.Default()

	r.POST("/ierror/message", func(c *gin.Context) {
		var msg gowon.Message

		if err := c.BindJSON(&msg); err != nil {
			log.Printf("Error: unable to bind message to json: %v", err)
			return
		}

		returnMsg := &gowon.Message{
			Module: moduleName,
			Msg:    "ierror errored",
			Dest:   msg.Dest,
		}

		c.IndentedJSON(http.StatusInternalServerError, returnMsg)
	})

	r.POST("/ierror2/message", func(c *gin.Context) {
		var msg gowon.Message

		if err := c.BindJSON(&msg); err != nil {
			log.Printf("Error: unable to bind message to json: %v", err)
			return
		}

		returnMsg := &gowon.Message{
			Module: moduleName,
			Msg:    "",
			Dest:   msg.Dest,
		}

		c.IndentedJSON(http.StatusInternalServerError, returnMsg)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
