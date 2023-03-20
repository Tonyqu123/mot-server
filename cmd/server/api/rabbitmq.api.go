package api

import (
	"fmt"
	"log"
	"net/url"
	"sync"
	"time"

	"github.com/tony/mot-server/cmd/server/config"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go"

	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQAPI struct{}

var (
	rabbitMQOnce sync.Once
	channel      *amqp.Channel
	queue        amqp.Queue
)

func InitRabbitMQOrDie() error {
	rabbitMQOnce.Do(
		func() {
			channel, queue = initRabbitMQOrDie()
		})
	return nil
}

func initRabbitMQOrDie() (*amqp.Channel, amqp.Queue) {
	rabbitMQInfo := config.GetRabbitMQFromEnv()
	// 连接到 RabbitMQ 服务器
	conn, err := amqp.Dial("amqp://" + rabbitMQInfo.User + ":" + rabbitMQInfo.Password + "@" + rabbitMQInfo.Endpoint + "/")
	if err != nil {
		log.Fatalf("FAIL TO intialize rabbitMQ, err: %s", err)
	}
	//defer conn.Close()

	// 创建一个channel来传递消息
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("FAIL TO intialize rabbitMQ channel, err: %s", err)
	}

	//defer ch.Close() // 函数 return 之后关闭 连接
	fmt.Println("initRabbitMQOrDie success")

	// 发送前，我们必须声明一个队列供我们发送，然后才能向队列发送消息
	// 设置与生产者相同，首先打开一个连接和一个Channel，并声明我们要消费的队列。请注意，这与发送的队列相匹配
	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		log.Fatalf("FAIL TO intialize rabbitMQ channel, err: %s", err)
	}

	return ch, q
}

func (a RabbitMQAPI) SendMessage(c *gin.Context) {
	fmt.Println("Hello World!")
	// 声明队列是幂等的——只有在它不存在的情况下才会创建它。消息内容是一个字节数组，因此你可以编写任何内容。
	body := "Hello World!"
	err := channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

	fmt.Println("if err != nil")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to publish a message", "data": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "RabbitMQ connect success"})
}

func (a RabbitMQAPI) ReceiveMessage(c *gin.Context) {
	// 取消息。确保使用消息之前队列已经存在。
	// 在goroutine中读取来自channel （由amqp :: Consume返回）的消息（持续不断地读取）。
	msgs, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Failed to declare a queue", "data": err.Error()})
		//failOnError(err, "Failed to register a consumer")
		return
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func (a RabbitMQAPI) GetMinio(c *gin.Context) {
	// Make a new bucket called mymusic.
	bucketName := "mymusic"
	location := "us-east-1"

	err := minioClient.MakeBucket(bucketName, location)
	if err != nil {
		// Check to see if we already own this bucket (which happens if you run this twice)
		exists, err := minioClient.BucketExists(bucketName)
		if err == nil && exists {
			log.Printf("We already own %s\n", bucketName)
		} else {
			log.Fatalln(err)
		}
	} else {
		log.Printf("Successfully created %s\n", bucketName)
	}

	// Upload the zip file
	objectName := "360210823142729440.png"
	filePath := "/Users/litingting/Desktop/360210823142729440.png"
	contentType := "application/png"

	// Upload the zip file with FPutObject
	n, err := minioClient.FPutObject(bucketName, objectName, filePath, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatalln(err)
	}

	reqParams := make(url.Values)
	expires := time.Second * 24 * 60 * 60
	presignedURL, err := minioClient.PresignedGetObject(bucketName, objectName, expires, reqParams)
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("下载链接：%s", presignedURL)

	log.Printf("Successfully uploaded %s of size %d\n", objectName, n)
}

//func failOnError(err error, msg string) {
//	if err != nil {
//		panic(err)
//	}
//}
