package route

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/riceChuang/jbtkLineBot/common"
	"github.com/riceChuang/jbtkLineBot/line"
	"github.com/riceChuang/jbtkLineBot/service"
	"net/http"
	"strings"
)

func CallbackHandler(w http.ResponseWriter, r *http.Request) {

	lineClient := line.GetBotClient()
	events, err := lineClient.ParseRequest(r)
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
				switch common.MessageKey(strings.ToLower(message.Text)) {
				case common.MessagePtt:
					service.ReplyPttMessage(event)
				case common.MessageDcard:
					service.ReplyDcardMessage(event)
				case common.MessageDcards:
					service.ReplyDcardMapMessage(event)
				case common.MessageJoke:
					service.ReplyJokeMessage(event)
				case common.MessagePorn:
					service.ReplyPornMessage(event)
				case common.MessageLong:
					service.ReplyMessageLen(event)
				default:
				}
			}
		}
	}
}
