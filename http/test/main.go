package main

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	pollerHttp "gitlab.com/ruangguru/polling/http"
)

func main() {
	payload, _ := json.Marshal(map[string]string{
		"janus":       "ping",
		"transaction": "123",
	})

	req, _ := http.NewRequest("POST", "http://localhost:7088/admin", bytes.NewReader(payload))

	agent := pollerHttp.Agent(req, 1, 1)

	ctx := context.Background()

	chResp := make(chan interface{})

	go agent.Run(ctx, chResp)

	for val := range chResp {
		value := val.(*pollerHttp.Response)
		log.Printf("\n\n%+v", value.Err)
		if value.Resp != nil {
			bodyBytes, err := ioutil.ReadAll(value.Resp.Body)
			if err != nil {
				log.Fatal(err)
			}
			bodyString := string(bodyBytes)
			log.Printf("\n\n%+v", bodyString)
		}
	}
}
