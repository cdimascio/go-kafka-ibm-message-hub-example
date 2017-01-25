package kafka

import (
	"crypto/tls"
	"fmt"
	"github.com/Shopify/sarama"
	"log"
	"os"
)

func init() {
	sarama.Logger = log.New(os.Stdout, "[sarama] ", log.LstdFlags)
}

type KafkaClient struct {
	client   sarama.Client
	producer sarama.SyncProducer
	consumer sarama.Consumer
}

func NewClient(user, pass, apiKey string, brokers []string) (KafkaClient, error) {
	emptyClient := KafkaClient{}
	tlsConfig, err := tlsConfig("server.pem", "server.key")
	if err != nil {
		return emptyClient, err
	}

	config := sarama.NewConfig()
	config.ClientID = apiKey
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	config.Net.TLS.Enable = true
	config.Net.TLS.Config = tlsConfig
	config.Net.SASL.User = user
	config.Net.SASL.Password = pass
	config.Net.SASL.Enable = true

	client, err := sarama.NewClient(brokers, config)
	if err != nil {
		return emptyClient, err
	}

	producer, err := sarama.NewSyncProducerFromClient(client)
	if err != nil {
		return emptyClient, err
	}

	consumer, err := sarama.NewConsumerFromClient(client)
	if err != nil {
		return emptyClient, err
	}

	return KafkaClient{client: client, producer: producer, consumer: consumer}, nil
}

func (client KafkaClient) SendMessage(topic, msg string) error {
	pMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(msg),
	}

	partition, offset, err := client.producer.SendMessage(pMsg)
	if err != nil {
		return err
	} else {
		fmt.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
		return nil
	}
}

func (client KafkaClient) ListenAndConsume(topic string, callback func(*sarama.ConsumerMessage)) error {
	partitionConsumer, err := client.consumer.ConsumePartition(topic, 0, sarama.OffsetNewest)
	if err != nil {
		return err
	}

	defer func() {
		if err := partitionConsumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	for {
		select {
		case msg := <-partitionConsumer.Messages():
			callback(msg)
		}
	}
}

func (client KafkaClient) Close() {
	if err := client.producer.Close(); err != nil {
		log.Fatalln(err)
	}
	if err := client.consumer.Close(); err != nil {
		log.Fatalln(err)
	}
	if err := client.client.Close(); err != nil {
		log.Fatalln(err)
	}
}

func tlsConfig(certFile, keyFile string) (*tls.Config, error) {
	cer, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}

	return &tls.Config{Certificates: []tls.Certificate{cer}}, nil
}
