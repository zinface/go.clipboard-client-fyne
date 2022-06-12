package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

type ClipBoard struct {
	Data     string    `json:"data,omitempty"`
	Mime     string    `json:"mime,omitempty"`
	CreateAt time.Time `json:"create_at"`
}

type ClipboardApi struct {
	Address string
}

func (api ClipboardApi) Get() (ClipBoard, error) {
	resp, err := http.Get(api.Address + "/clipboard")
	if err != nil {
		return ClipBoard{}, err
	}
	defer resp.Body.Close()

	var c ClipBoard
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return c, err
	}

	log.Println("api.Get():", c.Mime, c.CreateAt)

	return c, nil
}

func (api ClipboardApi) Set(content, mime string) (ClipBoard, error) {
	resp, err := http.Post(api.Address+"/clipboard", "application/json",
		strings.NewReader(`{"data":"`+content+`","mime":"`+mime+`"}`))
	if err != nil {
		return ClipBoard{}, err
	}
	defer resp.Body.Close()

	var c ClipBoard
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return c, err
	}

	if c.Mime != mime {
		return c, fmt.Errorf("mime type not match")
	}

	log.Println("api.Set():", c.Mime, c.CreateAt)

	return c, nil
}

func (api ClipboardApi) Info() (ClipBoard, error) {
	resp, err := http.Get(api.Address + "/clipboard/info")
	if err != nil {
		return ClipBoard{}, err
	}
	defer resp.Body.Close()

	var c ClipBoard
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return c, err
	}
	log.Println("api.Info():", c.Mime, c.CreateAt)

	return c, nil
}
