package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/bot-api/telegram"
	"github.com/orcaman/concurrent-map"
	"golang.org/x/net/context"
	"log"
	"net/http"
	"time"
)

type Orders struct {
	Ready     []string `json:"ready"`
	InProcess []string `json:"inProcess"`
}

var tr = &http.Transport{
	TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
}

var myClient = &http.Client{Timeout: 10 * time.Second, Transport: tr}

func getJson(url string, target interface{}) error {
	r, err := myClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

func main() {
	token := flag.String("token", "43", "telegram bot token")
	debug := flag.Bool("debug", false, "show debug information")
	flag.Parse()

	if *token == "" {
		log.Fatal("token flag required")
	}

	api := telegram.New(*token)
	api.Debug(*debug)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	m := cmap.New()

	ticker := time.NewTicker(10 * time.Second)
	go func() {
		for range ticker.C {
			newOrders := new(Orders)
			getJson("http://kotikicanteen.ru/orders", newOrders)

			fmt.Printf("%+v\n", newOrders)

			for _, r := range newOrders.Ready {
				if messange, ok := m.Get(r); ok {
					msger := messange.(*telegram.Message)
					fmt.Printf("%+v\n", msger)

					newMesage := &telegram.MessageCfg{
						BaseMessage: telegram.BaseMessage{
							BaseChat: telegram.BaseChat{
								ID: msger.Chat.ID,
							},
						},
						Text: "Ready: " + r,
					}

					if _, err := api.Send(ctx, newMesage); err != nil {
						log.Printf("send error: %v", err)
					}

					m.Remove(r)
				}
			}
		}
	}()

	if user, err := api.GetMe(ctx); err != nil {
		log.Panic(err)
	} else {
		log.Printf("bot info: %#v", user)
	}

	updatesCh := make(chan telegram.Update)

	go telegram.GetUpdates(ctx, api, telegram.UpdateCfg{
		Timeout: 10, // Timeout in seconds for long polling.
		Offset:  0,  // Start with the oldest update
	}, updatesCh)

	for update := range updatesCh {
		log.Printf("got update from %s", update.Message.From.Username)
		if update.Message == nil {
			continue
		}
		if len(update.Message.Text) > 4 {
			log.Println("long message")
			continue
		}

		m.Set(update.Message.Text, update.Message)
	}
}
