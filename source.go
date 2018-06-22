package main

import (
	"encoding/json"
	"fmt"
	"github.com/zhouhui8915/go-socket.io-client"
	"io/ioutil"
	"log"
	"net/http"
)

type JoinChannelMessage struct {
	Name string `json:"name"`
}
type LoginMessage struct {
	Name     string `json:"name"`
	Password string `json:"pw"`
}
type chatMessage struct {
	Text string `json:"msg"`
}
type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func main() {
	channel := "fullmoviesonyoutube"
	url, _ := cytubeUrl(channel)

	opts := &socketio_client.Options{
		Transport: "websocket",
		Query:     make(map[string]string),
	}

	client, err := socketio_client.NewClient(url, opts)
	if err != nil {
		log.Printf("NewClient error:%v\n", err)
		return
	}
	client.On("chatMsg", func() {
		log.Printf("message")
	})
	client.On("disconnection", func() {
		log.Printf("on disconnect\n")
	})
	client.On("event", func(msg string) {
		log.Printf("on event\n", msg)
	})
	client.On("chatMsg", func(msg string) {
		log.Printf("on message:%v\n", msg)
	})
	var joinChannel JoinChannelMessage
	var login LoginMessage
	var chat chatMessage
	chat.Text = "test"
	login.Name = "pookeytesting1"
	joinChannel.Name = channel
	//client.Emit("initChannelCallbacks")
	client.Emit("joinChannel", joinChannel)
	client.Emit("login", login)

	for {
		client.Emit("chatMsg", chat)
	}
}

func cytubeUrl(channel string) (string, error) {
	data := struct {
		Items []struct {
			Url    string `json:"url"`
			Secure bool   `json:"secure"`
		} `json:"servers"`
	}{}
	url := fmt.Sprintf("https://cytu.be/socketconfig/%s.json", channel)
	getJson(url, &data)
	for _, server := range data.Items {
		if !server.Secure {
			return server.Url, nil
		}
	}
	return "", nil
}

func getJson(url string, target interface{}) error {
	httpClient := http.Client{}
	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &target)
	if err != nil {
		return err
	}

	return nil
}
