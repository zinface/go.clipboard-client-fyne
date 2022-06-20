package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

type Clipboard struct {
	Data     string    `json:"data,omitempty"`
	Mime     string    `json:"mime,omitempty"`
	CreateAt time.Time `json:"create_at"`
}

func (c Clipboard) String() string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "[%s] (%s)", c.Mime, c.CreateAt.Format("2006-01-02 15:04:05"))
	if len(c.Data) > 0 {
		fmt.Fprintf(&sb, " (%d bytes)", len(c.Data))
		switch c.Mime {
		case "text/plain":
			fmt.Fprintf(&sb, " => %q", c.Data)
		case "image/png":
			// fmt.Fprintf(&sb, " => (%d x %d)", ...)
		}
	}
	return sb.String()
}

type ClipboardAPI struct {
	Address string
}

func (api ClipboardAPI) Get() (Clipboard, error) {
	resp, err := http.Get(api.Address + "/clipboard")
	if err != nil {
		return Clipboard{}, err
	}
	defer resp.Body.Close()

	var c Clipboard
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return c, err
	}

	log.Debugf("api.Get(): %v", c)

	return c, nil
}

func (api ClipboardAPI) Set(content, mime string) (Clipboard, error) {
	resp, err := http.Post(api.Address+"/clipboard", "application/json",
		strings.NewReader(`{"data":"`+content+`","mime":"`+mime+`"}`))
	if err != nil {
		return Clipboard{}, err
	}
	defer resp.Body.Close()

	var c Clipboard
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return c, err
	}

	if c.Mime != mime {
		return c, fmt.Errorf("mime type not match")
	}

	log.Debugf("api.Set(): %v", c)

	return c, nil
}

func (api ClipboardAPI) Info() (Clipboard, error) {
	resp, err := http.Get(api.Address + "/clipboard/info")
	if err != nil {
		return Clipboard{}, err
	}
	defer resp.Body.Close()

	var c Clipboard
	if err := json.NewDecoder(resp.Body).Decode(&c); err != nil {
		return c, err
	}
	log.Tracef("api.Info(): %v", c)

	return c, nil
}
