package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

//structure of the response from the Real time messagin api 
type rmtStartResponse struct {
	Ok    bool    `json:"ok"`
	Url   string  `json:"url"`
	Error string  `json:"error"`
	Self  rmtSelf `json:"self"`
}

type rmtSelf struct {
	Id string `json:"id"`
}

//structure of message sent by the user
type Message struct {
	Id      uint64 `json:"id"`
	Type    string `json:"type"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

//Get details necessary to start a websocket connection for real time messaging
func getConnectionDetails(token string) (webSocketUrl string, Id string, err error) {

	slackStartRmtUrl := fmt.Sprintf("https://slack.com/api/rtm.start?token=%s", token)
	resp, err := http.Get(slackStartRmtUrl)

	if err != nil {
		fmt.Println("An error has occured", err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	var rmtResponse rmtStartResponse
	err = json.Unmarshal(body, &rmtResponse)

	if err != nil {
		return
	}

	if !rmtResponse.Ok {
		err = fmt.Errorf("Error from slack API", rmtResponse.Error)
		return
	}

	webSocketUrl = rmtResponse.Url
	Id = rmtResponse.Self.Id
	return
}

//Create a websocket connection to receive and send messages
func initiateConnection(token string) (webSocket *websocket.Conn, id string) {
	webSocketUrl, id, err := getConnectionDetails(token)
	if err != nil {
		log.Fatal(err)
	}

	webSocket, err = websocket.Dial(webSocketUrl, "", "https://api.slack.com")

	if err != nil {
		log.Fatal(err)
	}

	return
}

//receive messages via the websocket
func receiveMessage(webSocket *websocket.Conn) (msg Message, err error) {
	err = websocket.JSON.Receive(webSocket, &msg)

	return
}

//send messages via the websocket
func sendMessage(webSocket *websocket.Conn, msg Message) (err error) {
	err = websocket.JSON.Send(webSocket, &msg)
	return

}
