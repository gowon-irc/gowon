package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gowon-irc/go-gowon"
)

const (
	moduleName = "module3"
	moduleHelp = "return message in cyan"
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
			Msg:    fmt.Sprintf("{cyan}%s{clear}", msg.Args),
			Dest:   msg.Dest,
		}

		c.IndentedJSON(http.StatusOK, returnMsg)
	})

	r.GET("/help", func(c *gin.Context) {
		c.IndentedJSON(http.StatusOK, &gowon.Message{
			Module: moduleName,
			Msg:    moduleHelp,
		})
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
