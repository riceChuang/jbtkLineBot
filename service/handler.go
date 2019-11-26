package service

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/riceChuang/jbtkLineBot/crawler"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var bot *linebot.Client

func InitBotClient() {
	var err error
	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
}

func CallbackHandler(w http.ResponseWriter, r *http.Request) {

	events, err := bot.ParseRequest(r)
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
				switch strings.ToLower(message.Text) {
				case MessagePtt:
					ReplyPttMessage(event)
				case MessageDcard:
					ReplyDcardMessage(event)
				case MessageDcards:
					ReplyDcardMapMessage(event)
				case MessageJoke:
					ReplyJokeMessage(event)
				case MessagePorn:
					ReplyPornMessage(event)
				case MessageLong:
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("beauty len :"+strconv.Itoa(crawler.ImageLength)+"decard len :"+strconv.Itoa(crawler.DcardImageLengh)+"joker len :"+strconv.Itoa(crawler.JokerLenght))).Do(); err != nil {
						log.Print(err)
					}
				default:
				}
			}
		}
	}
}
