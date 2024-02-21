package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type returnMsg struct {
	Msg string `json:"msg"`
}

func main() {
	r := gin.Default()
	r.GET("/message", func(c *gin.Context) {
		input := c.Query("msg")

		log.Println(input)

		output := returnMsg{
			Msg: fmt.Sprintf("{cyan}%s{clear}", input),
		}

		c.JSON(http.StatusOK, output)
	})

	r.Run(":8080")
}
