package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/line/line-bot-sdk-go/linebot"
)

var bot *linebot.Client

func main() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/back", callbackHandler)
	port := os.Getenv("PORT")
	addr := ""
	if port != ""{
		addr = fmt.Sprintf(":%s", port)
	}else{
		addr = ":3000"
	}

	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	events, err := bot.ParseRequest(r)
	//config := types.New("app.yml")
	//fmt.Println(config.Images[0])
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			w.WriteHeader(400)
		} else {
			w.WriteHeader(500)
		}
		return
	}

	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				if message.Text == "抽" {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage("https://imgur.com/JPfUBeF", "https://imgur.com/JPfUBeF")).Do(); err != nil {
						log.Print(err)
					}
				} else if message.Text == "機吧毛" {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage("https://i.imgur.com/khCKl58.jpg", "https://i.imgur.com/khCKl58.jpg")).Do(); err != nil {
						log.Print(err)
					}
				}

			}
		}
	}
}
