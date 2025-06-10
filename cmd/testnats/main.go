package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

func main() {
	ctx := gctx.New()
	nc, err := nats.Connect("nats://nats-server:4222")
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	js, err := jetstream.New(nc)
	if err != nil {
		log.Fatal(err)
	}
	js.CreateStream(ctx, jetstream.StreamConfig{
		Name:       "STREAM_NAME",
		Subjects:   []string{"subjects"},
		Storage:    jetstream.FileStorage, // 默认文件存储
		MaxMsgSize: 10 * 1024 * 1024,
		Retention:  jetstream.InterestPolicy, // 如果有多个消费者订阅了相同的主题，等每个消费者都消费确认后删除消息
	})
	// 测试消费者列表权限
	_, err = js.Consumer(ctx, "STREAM_NAME", "CONSUMER_NAME")
	if err != nil {
		fmt.Printf("Consumer list permission denied: %v\n", err)
	} else {
		fmt.Println("Has consumer list permission")
	}

	// 测试消费者创建权限
	_, err = js.CreateOrUpdateConsumer(ctx, "STREAM_NAME", jetstream.ConsumerConfig{
		Durable:       "TEST_PERM_CHECK_DURABLE",
		AckPolicy:     jetstream.AckExplicitPolicy,
		Name:          "TEST_PERM_CHECK_DURABLE",
		FilterSubject: "subjects",
		DeliverPolicy: jetstream.DeliverNewPolicy,
		AckWait:       30 * time.Second,
	})
	if err != nil {
		fmt.Printf("Consumer create permission denied: %v\n", err)
	} else {
		fmt.Println("Has consumer create permission")
	}
	consumer, err := js.CreateOrUpdateConsumer(ctx, "STREAM_NAME", jetstream.ConsumerConfig{
		Durable:       "TEST_PERM_CHECK_DURABLE222",
		AckPolicy:     jetstream.AckExplicitPolicy,
		Name:          "TEST_PERM_CHECK_DURABLE222",
		FilterSubject: "subjects",
		DeliverPolicy: jetstream.DeliverNewPolicy,
		AckWait:       30 * time.Second,
	})
	if err != nil {
		fmt.Printf("Consumer create permission denied: %v\n", err)
	} else {
		fmt.Println("Has consumer create permission")
	}
	iter, err := consumer.Messages(jetstream.PullMaxMessages(1))
	if err != nil {
		fmt.Println(err)
	}
	go func() {
		for {
			msg, err := iter.Next()
			if err != nil {
				fmt.Println("nets err:", err)
				break
			}
			msg.Ack()
		}
	}()
	err = js.DeleteConsumer(ctx, "STREAM_NAME", "TEST_PERM_CHECK_DURABLE")
	if err != nil {
		fmt.Printf("Consumer delete permission denied: %v\n", err)
	} else {
		fmt.Println("Has consumer delete permission")
	}
	err = js.DeleteConsumer(ctx, "STREAM_NAME", "TEST_PERM_CHECK_DURABLE222")
	if err != nil {
		fmt.Printf("Consumer delete permission denied: %v\n", err)
	} else {
		fmt.Println("Has consumer delete permission 222")
	}
}
