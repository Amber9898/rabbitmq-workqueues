package mqUtils

import (
	"context"
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"sync"
	"time"
)

const (
	QUEUE_NAME = "hello"
	MQ_ADDR    = "192.168.111.3"
	MQ_USER    = "admin"
	MQ_PWD     = "admin"
)

type MQPublisherManager struct {
	Con     *amqp.Connection
	Ch      *amqp.Channel
	Queue   *amqp.Queue
	Durable bool
}

var mqPublisherInstance *MQPublisherManager
var publisherOnce sync.Once

func GetMQPublisherInstance() *MQPublisherManager {
	return mqPublisherInstance
}

// InitPublisher durable: 是否持久化队列 和 消息
func InitPublisher(durable bool) {
	publisherOnce.Do(func() {
		initPublisher(durable)
	})
}

func initPublisher(durable bool) {
	mqPublisherInstance = &MQPublisherManager{
		Con:     nil,
		Ch:      nil,
		Queue:   nil,
		Durable: durable,
	}
	//1. 建立mq链接
	url := fmt.Sprintf("amqp://%s:%s@%s:5672/", MQ_USER, MQ_PWD, MQ_ADDR)
	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Println("mq connection error: ", err)
	}
	mqPublisherInstance.Con = conn

	//2. 创建信道
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("channel creation error: ", err)
	}
	mqPublisherInstance.Ch = ch

	//3. 声明队列
	q, err := ch.QueueDeclare(
		QUEUE_NAME, //队列名
		durable,    //是否持久化到磁盘
		false,      //当最后一个消费者断开连接后，是否自动删除
		false,      //是否只允许一个消费者消费
		false,      //
		nil)
	if err != nil {
		fmt.Println("queue declaration error: ", err)
	}
	mqPublisherInstance.Queue = &q
}

func (mq *MQPublisherManager) Publish(msg string) {
	context, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	body := msg

	mode := amqp.Transient //消息模式
	if mq.Durable {
		mode = amqp.Persistent
	}

	err := mq.Ch.PublishWithContext(
		context,
		"",
		mq.Queue.Name,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: mode,
			ContentType:  "text/plain",
			Body:         []byte(body),
		},
	)
	if err != nil {
		fmt.Println("queue publishing error: ", err)
		return
	}
	fmt.Println("finish publishing mq messages...")
}
