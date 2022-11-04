package main

import (
	"awesomeProject/common/mqUtils"
	"fmt"
)

func main() {
	fmt.Println("Please input something...")
	fmt.Println("enter break to stop the  loop")
	//初始化
	durable := true
	mqUtils.InitPublisher(durable)
	publisher := mqUtils.GetMQPublisherInstance()
	for {
		var val string
		fmt.Scanf("%s", &val)
		if val == "break" {
			fmt.Println("stop the loop")
			break
		} else if val == "\\n" || val == "\\t" || val == "" {
			continue
		} else {
			fmt.Println("Please input something...")
			fmt.Println("enter break or press ctrl+c to stop the  loop")
			fmt.Println("--------->", val)
			go publisher.Publish(val)
		}
	}
}
