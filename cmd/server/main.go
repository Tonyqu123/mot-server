package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"
	"github.com/tony/mot-server/cmd/server/middleware"
	"github.com/tony/mot-server/cmd/server/router"

	"github.com/gin-gonic/gin"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	r := gin.Default()
	r.Use(middleware.Cors()) //开启中间件 允许使用跨域请求

	r.POST("/upload", func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")

		// Source
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, "get form err: %s", err.Error())
			return
		}

		filename := filepath.Base(file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			c.String(http.StatusBadRequest, "upload file err: %s", err.Error())
			return
		}

		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		if err != nil {
			c.String(http.StatusBadRequest, "connect err: %s", err.Error())
			return

		}
		defer conn.Close()
		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		q, err := ch.QueueDeclare(
			"hello", // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		failOnError(err, "Failed to declare a queue")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		body := "Hello moter, a vedio uploaded: %s"
		err = ch.PublishWithContext(ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(fmt.Sprintf(body, filename)),
			})
		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s\n", body)
		c.String(http.StatusOK, "File %s uploaded successfully with fields name=%s and email=%s.", file.Filename, name, email)
	})
	r2 := router.Router{}
	r2.RegisterAPI(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
