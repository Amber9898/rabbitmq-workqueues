package mqUtils

import (
	"fmt"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func StartConsumer(durable bool) {
	//1. 建立mq链接
	url := fmt.Sprintf("amqp://%s:%s@%s:5672/", MQ_USER, MQ_PWD, MQ_ADDR)
	conn, err := amqp.Dial(url)
	if err != nil {
		fmt.Println("mq connection error: ", err)
		return
	}
	defer conn.Close()

	//2. 创建信道
	ch, err := conn.Channel()
	if err != nil {
		fmt.Println("channel creation error: ", err)
		return
	}
	defer ch.Close()

	//3. 声明队列
	//为什么消费者也声明队列：
	//Because we might start the consumer before the publisher,
	//we want to make sure the queue exists before we try to
	//consume messages from it.
	q, err := ch.QueueDeclare(
		QUEUE_NAME, //队列名
		durable,    //是否持久化到磁盘
		false,      //当最后一个消费者断开连接后，是否自动删除
		false,      //是否只允许一个消费者消费
		false,      //
		nil)
	if err != nil {
		fmt.Println("queue declaration error: ", err)
		return
	}

	//4. 消费队列消息
	msgCh, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil)
	if err != nil {
		fmt.Println("consume message error: ", err)
		return
	}
	fmt.Println("start a consumer....")
	forever := make(chan bool)
	go msgChannel(msgCh)
	<-forever

}

func msgChannel(ch <-chan amqp.Delivery) {

	for msg := range ch {
		fmt.Println("receive message from mq--->", string(msg.Body), msg.ContentType)
	}
	time.Sleep(100 * time.Millisecond)

	fmt.Println("end...")
}
