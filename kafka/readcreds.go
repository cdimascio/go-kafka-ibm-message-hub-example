package kafka

import (
	"encoding/json"
	"io/ioutil"
)

type MessageHubCredentials struct {
	Username     string   `json:"user"`
	Password     string   `json:"password"`
	KafkaBrokers []string `json:"kafka_brokers_sasl"`
	ApiKey       string   `json:"api_key"`
}

func MessageHubCreds() (MessageHubCredentials, error) {
	creds := MessageHubCredentials{}
	bytes, err := ioutil.ReadFile("message-hub-creds.json")
	if err != nil {
		return creds, err
	}

	err = json.Unmarshal(bytes, &creds)
	if err != nil {
		return creds, err
	}

	return creds, nil
}
