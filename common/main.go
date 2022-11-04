package main

import "awesomeProject/common/mqUtils"

func main() {
	forever := make(chan bool)
	durable := true
	go mqUtils.StartConsumer(durable)
	<-forever
}
