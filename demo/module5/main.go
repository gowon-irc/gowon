package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

type returnMsg struct {
	Msg string `json:"msg"`
}

func main() {
	r := gin.Default()
	r.GET("/message", func(c *gin.Context) {
		input := c.Query("msg")

		log.Println(input)

		output := returnMsg{
			Msg: reverse(input),
		}

		c.JSON(http.StatusOK, output)
	})

	r.Run(":8080")
}
