package main

import (
	"fmt"
	"flag"
	"github.com/Shopify/sarama"
	"log"
	"os"
	"strconv"
	"github.com/cdimascio/go-kakfa-ibm-message-hub-example/kafka"
)

func main() {
	clientType := flag.String("type", "consumer", "producer|consumer")
	topic := flag.String("topic", "mytopic", "mytopic")
	num := flag.Int("num", 10, "number of messages to produce. no op is type is consumer")
	flag.Parse()

	creds, err := kafka.MessageHubCreds()
	exitOnErr(err)

	client, err := kafka.NewClient(creds.Username, creds.Password, creds.ApiKey, creds.KafkaBrokers)
	exitOnErr(err)

	switch *clientType {
	case "producer":
		for i := 0; i < *num; i++ {
			err = client.SendMessage(*topic, "message "+strconv.Itoa(i))
		}
	case "consumer":
		err = client.ListenAndConsume(*topic, func(msg *sarama.ConsumerMessage) {
			log.Printf("Consumed message offset %d\n", msg.Offset)
		})
	}
	exitOnErr(err)
	fmt.Scanln()
	client.Close()
}

func exitOnErr(err error) {
	if err != nil {
		fmt.Printf("%s", err.Error())
		os.Exit(1)
	}
}
