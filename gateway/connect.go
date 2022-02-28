package gateway

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	Config "micrach/config"
)

// Make http request to gateway to tell the board id and description
func Connect() {
	requestBody, _ := json.Marshal(map[string]string{
		"id":   Config.App.Gateway.BoardId,
		"name": Config.App.Gateway.BoardDescription,
		"url":  Config.App.Gateway.Url,
	})
	url := Config.App.Gateway.Url + "/api/boards"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		log.Panicln(err)
	}
	req.Header.Set("Authorization", Config.App.Gateway.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Panicln(err)
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Panicln(err)
	}
	log.Println(string(body))
	log.Println("gateway - online")
}
