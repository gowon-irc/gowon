package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gowon-irc/go-gowon"
)

const (
	moduleName = "module2"
)

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

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
			Msg:    reverse(msg.Args),
			Dest:   msg.Dest,
		}

		c.IndentedJSON(http.StatusOK, returnMsg)
	})

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
