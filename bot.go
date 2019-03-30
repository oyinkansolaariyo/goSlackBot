package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var token = flag.String("token", "", "API token for the bot, check your dashboard")

func main() {
	//get token from flags passed in from the command line
	flag.Parse()
	if *token == "" {
		fmt.Fprintf(os.Stderr, "You have to provide a token")
		os.Exit(1)
	}
	
	//initialize a websocket connection
	webSocket, _ := initiateConnection(*token)

	for {
	//receive message and do anything with the message
		msg, err := receiveMessage(webSocket)

		if err != nil {
			log.Fatal(err)
		}

		if msg.Type == "message" {

			fmt.Println("Message type", msg.Type)
			fmt.Println("Message text", msg.Text)
			fmt.Println("Message channel", msg.Channel)

			sendMessage(webSocket, msg)
		}

	}

}
