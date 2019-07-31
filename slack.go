package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	POST_URL = "https://slack.com/api/chat.postMessage"
)

type EventURLVerifyReqBody struct {
	Token     string `json:"token"`
	Challenge string `json:"challenge"`
	Type      string `json:"type"`
}

type EventBody struct {
	Token string          `json:"token"`
	Event json.RawMessage `json:"event"`
	Type  string          `json:"type"`
}
type Event struct {
	Type     string          `json:"type"`
	Event_ts string          `json:"event_ts"`
	User     string          `json:"user"`
	Bot_id   string          `json:"bot_id"`
	Text     string          `json:"text"`
	Channel  string          `json:"channel"`
	Files    json.RawMessage `json:"files"`
}

type File struct {
	Name string `json:"name"`
	Size int    `json:"size"`
	URL  string `json:"url_private"`
}

func postMessage(channel, text string) error {
	values := url.Values{}
	values.Set("text", text)
	values.Add("channel", channel)
	values.Add("token", TOKEN)

	req, err := http.NewRequest("POST", POST_URL, strings.NewReader((values.Encode())))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	cli := &http.Client{}
	res, err := cli.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	return nil
}

var TOKEN string

func main() {
	TOKEN = os.Getenv("SLACK_TOKEN")
	/*
		http.HandleFunc("/urlverify", func(w http.ResponseWriter, r *http.Request) {
			if r.Method != "POST" {
				w.WriteHeader(400)
				return
			}
			fmt.Println(r)
			buff := new(bytes.Buffer)
			io.Copy(buff, r.Body)
			fmt.Println(buff)
			body := new(EventURLVerifyReqBody)
			if err := json.Unmarshal(buff.Bytes(), &body); err != nil {
				fmt.Println(err)
			}
			fmt.Println(body)
			w.WriteHeader(200)
			w.Write([]byte(body.Challenge))
		})
	*/

	http.HandleFunc("/urlverify", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(400)
			return
		}

		buff := new(bytes.Buffer)
		io.Copy(buff, r.Body)

		fmt.Println(string(buff.Bytes()))

		eventReq := new(EventBody)
		if err := json.Unmarshal(buff.Bytes(), &eventReq); err != nil {
			fmt.Println(err)
		}
		fmt.Print("EventReq: ")
		fmt.Println(eventReq.Event)
		fmt.Println(string(eventReq.Event))

		if eventReq.Type == "url_verification" {
			body := new(EventURLVerifyReqBody)
			if err := json.Unmarshal(buff.Bytes(), &body); err != nil {
				fmt.Println(err)
			}
			fmt.Println(body)
			w.WriteHeader(200)
			w.Write([]byte(body.Challenge))
		}

		event := new(Event)
		if err := json.Unmarshal(eventReq.Event, &event); err != nil {
			fmt.Println(err)
		}
		fmt.Print("Event: ")
		fmt.Print(event.Type)
		fmt.Println(event.Event_ts)
		fmt.Println(event.User)
		fmt.Println(event.Text)
		fmt.Println(event.Channel)
		fmt.Println(event)

		fmt.Println(string(event.Files))

		files := new([]File)

		// if Files attached
		if len(event.Files) > 0 {
			if err := json.Unmarshal(event.Files, files); err != nil {
				fmt.Println(err)
			}

			for _, i := range *files {
				fmt.Println(i)
			}
		}

		fmt.Println(event.Bot_id)
		if event.Bot_id == "" {
			fmt.Println(postMessage(event.Channel, "<@"+event.User+">"+event.Text))
		}

		w.WriteHeader(200)
	})
	http.ListenAndServe(":8080", nil)
}
