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


	// 测试连接 MQ，并创建 Channel，往里面发消息
	r.GET("/test-mq", func(c *gin.Context) {
		// 连接到 RabbitMQ 服务器
		conn, err := amqp.Dial("amqp://admin:123456@localhost:5672/")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "connect err", "data": err.Error()})
		}
		defer conn.Close()

		// 创建一个channel来传递消息
		ch, err := conn.Channel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to open a channel", "data": err.Error()})
		}
		defer ch.Close()

		// 发送前，我们必须声明一个队列供我们发送，然后才能向队列发送消息
		q, err := ch.QueueDeclare(
			"hello", // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to declare a queue", "data": err.Error()})
		}
		
		// 声明队列是幂等的——只有在它不存在的情况下才会创建它。消息内容是一个字节数组，因此你可以编写任何内容。
		body := "Hello World!"
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing {
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to publish a message", "data": err.Error()})
		}

		c.JSON(http.StatusOK, gin.H{"message": "RabbitMQ connect success"})
	})


	// 测试接受 MQ 中的消息
	r.GET("receive-mq-msg", func(c *gin.Context) {
		// 设置与生产者相同，首先打开一个连接和一个Channel，并声明我们要消费的队列。请注意，这与发送的队列相匹配
		conn, err := amqp.Dial("amqp://admin:123456@localhost:5672/")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to connect to RabbitMQ", "data": err.Error()})
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to open a channel", "data": err.Error()})
		}
		defer ch.Close()

		q, err := ch.QueueDeclare(
			"hello", // name
			false,   // durable
			false,   // delete when usused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to declare a queue", "data": err.Error()})
		}

		// 取消息。确保使用消息之前队列已经存在。
		// 在goroutine中读取来自channel （由amqp :: Consume返回）的消息（持续不断地读取）。
		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)
		failOnError(err, "Failed to register a consumer")
	
		var forever chan struct{}
	
		go func() {
			for d := range msgs {
				log.Printf("Received a message: %s", d.Body)
			}
		}()
	
		log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
		<-forever
	})

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

func failOnError(err error, msg string) {
	if err != nil {
		panic(err)
	}
}
