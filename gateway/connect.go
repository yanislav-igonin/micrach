package gateway

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"micrach/config"
)

func Connect() {
	body, _ := json.Marshal(map[string]string{
		"name": config.App.Gateway.BoardName,
		"url":  config.App.Gateway.BoardUrl,
	})
	req, err := http.NewRequest(http.MethodPost, config.App.Gateway.Url+"/api/boards/me", bytes.NewReader(body))
	if err != nil {
		log.Println("gateway - connect error:", err)
		return
	}
	req.Header.Set("Authorization", "Bearer "+config.App.Gateway.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println("gateway - connect error:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		log.Println("gateway - connect failed:", resp.StatusCode, string(b))
		return
	}
	log.Println("gateway - online")
}
