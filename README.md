# go-kafka-ibm-message-hub-example

This sample show how to use IBM Message Hub with GoLang using the sarama kafka client.

## Setup and Run

1. Create an [IBM Message Hub](https://console.ng.bluemix.net/catalog/services/message-hub) Service instance via IBM [Bluemix](https://www.bluemix.net[]())

2. Clone this repo

	- `git clone https://github.ibm.com/dimascio/go-kafka-ibm-message-hub-example`

3. Create your own certificate

	Generated private key:

	`openssl genrsa -out server.key 2048`

	To generate a certificate:

	`openssl req -new -x509 -key server.key -out server.pem -days 3650`

	This will create `server.pem` and a `server.key`

4. `cd <project-root>` where `<project-root>` is where you cloned this repo
5. Copy `server.pem` and a `server.key` to `<project-root>`
6. Move `message-hub-creds.json.sample` to `message-hub-creds.json`
7. Populate `message-hub-creds.json` with credentials from step 1.
4. Run `go build`
5. Start the consumer.

 	`go-kafka-ibm-message-hub-example -type=consumer -topic=mytopic`

6. Start the producer

	`go-kafka-ibm-message-hub-example -type=producer -topic=mytopic -num=20`
