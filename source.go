package main

import (
	"encoding/json"
	"fmt"
	"github.com/graarh/golang-socketio"
	"github.com/graarh/golang-socketio/transport"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Message struct {
	Id      int    `json:"id"`
	Channel string `json:"channel"`
	Text    string `json:"text"`
}

func main() {
	url, err := cytube("fullmoviesonyoutube")
	if err != nil {
	}
	c, err := gosocketio.Dial(url, transport.GetDefaultWebsocketTransport())
	if err != nil {
		log.Fatal(err)
	}

	err = c.On("/message", func(h *gosocketio.Channel, args Message) {
		log.Println("--- Got chat message: ", args)
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnDisconnection, func(h *gosocketio.Channel) {
		log.Fatal("Disconnected")
	})
	if err != nil {
		log.Fatal(err)
	}

	err = c.On(gosocketio.OnConnection, func(h *gosocketio.Channel) {
		log.Println("Connected")
	})
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(1 * time.Second)
	time.Sleep(60 * time.Second)
	c.Close()

	log.Println(" [x] Complete")
}

func cytube(channel string) (string, error) {
	data := struct {
		Items []struct {
			Url    string `json:"url"`
			Secure bool   `json:"secure"`
		} `json:"servers"`
	}{}
	url := fmt.Sprintf("https://cytu.be/socketconfig/%s.json", channel)
	getJson(url, &data)
	for _, server := range data.Items {
		if server.Secure {
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
