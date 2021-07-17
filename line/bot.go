package line

import (
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/sirupsen/logrus"
	"os"
	"sync"
)

var (
	BotClient *linebot.Client
	botOnce   = sync.Once{}
)

func InitBotClient() {
	botOnce.Do(func() {
		var err error
		BotClient, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
		if err != nil {
			//logrus.Panicf("init linebot error:%v", err)
			return
		}
		logrus.Infof("Bot:", BotClient, " err:", err)
	})
}

func GetBotClient() *linebot.Client {
	return BotClient
}
