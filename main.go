package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/riceChuang/jbtkLineBot/boltdb"
	"github.com/riceChuang/jbtkLineBot/crawler"
	"github.com/riceChuang/jbtkLineBot/types"
	"math/rand"
	"strings"
)

var bot *linebot.Client

func main() {

	boltdb.Initialize()
	crawler.Initialize()

	config := types.New("./app.yml")
	var err error

	crawlerTypesMap := map[crawler.Type]string{
		crawler.Beauty:   config.BeautyUrl,
		crawler.DcardSex: config.DcardUrl,
		crawler.Joker:    config.JokerUrl,
	}

	for crawlerType := range crawlerTypesMap {
		crawlerWorker, err := crawler.GetCrawlerByType(crawlerType)
		if err != nil {
			fmt.Println(err)
		}
		go crawlerWorker.RunImage(crawlerTypesMap[crawlerType])
	}

	bot, err = linebot.New(os.Getenv("ChannelSecret"), os.Getenv("ChannelAccessToken"))
	log.Println("Bot:", bot, " err:", err)
	http.HandleFunc("/back", callbackHandler)
	port := os.Getenv("PORT")
	addr := ""
	if port != "" {
		addr = fmt.Sprintf(":%s", port)
	} else {
		addr = ":3000"
	}

	http.ListenAndServe(addr, nil)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
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
				if message.Text == "抽" {
					imageIndex := rand.Intn(crawler.ImageLengh)
					db := boltdb.DB()
					dbkey := fmt.Sprintf("beauty-%d", imageIndex)
					url := db.Read(dbkey)
					fmt.Printf("my image link: %v", url)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(url, url)).Do(); err != nil {
						log.Print(err)
					}
				} else if strings.ToLower(message.Text) == "d" {
					imageIndex := rand.Intn(crawler.DcardImageLengh)
					db := boltdb.DB()
					dbkey := fmt.Sprintf("dcard-%d", imageIndex)
					url := db.Read(dbkey)
					fmt.Printf("my image link: %v", url)
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(url, url)).Do(); err != nil {
						log.Print(err)
					}
				}else if message.Text == "笑"{
					imageIndex := rand.Intn(crawler.JokerLenght)
					dbkey := fmt.Sprintf("joker-%d", imageIndex)
					content := crawler.JokerMap[dbkey]
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(content)).Do(); err != nil {
						log.Print(err)
					}
				} else if message.Text == "長度" {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("beauty len :"+strconv.Itoa(crawler.ImageLengh)+"decard len :"+strconv.Itoa(crawler.DcardImageLengh)+"joker len :"+strconv.Itoa(crawler.JokerLenght))).Do(); err != nil {
						log.Print(err)
					}
				}
			}
		}
	}
}
