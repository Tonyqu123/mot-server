package api

import (
  "log"
	"github.com/gin-gonic/gin"
  // "github.com/streadway/amqp"
	amqp "github.com/rabbitmq/amqp091-go"
	"net/http"
)

type RabbitMQAPI struct {}

func (a RabbitMQAPI) SendMessage(c *gin.Context) {
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
}

func (a RabbitMQAPI) ReceiveMessage(c *gin.Context) {
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
}

func failOnError(err error, msg string) {
	if err != nil {
		panic(err)
	}
}