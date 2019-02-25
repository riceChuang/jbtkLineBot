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
	beautyCrawler, err := crawler.GetCrawlerByType(crawler.Beauty)
	if err != nil {
		fmt.Println(err)
	}
	go beautyCrawler.RunImage(config.BeautyUrl)

	dcardCrawler, err := crawler.GetCrawlerByType(crawler.DcardSex)
	if err != nil {
		fmt.Println(err)
	}
	go dcardCrawler.RunImage(config.DcardUrl)

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
				} else if message.Text == "機吧毛" {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage("https://i.imgur.com/khCKl58.jpg", "https://i.imgur.com/khCKl58.jpg")).Do(); err != nil {
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
				//} else if strings.ToLower(message.Text) == "ddd" {
				//	message := []*linebot.ImageMessage{}
				//	for i := 0; i < 3; i++ {
				//		imageIndex := rand.Intn(crawler.DcardImageLengh)
				//		db := boltdb.DB()
				//		dbkey := fmt.Sprintf("dcard-%d", imageIndex)
				//		url := db.Read(dbkey)
				//		fmt.Printf("my image link: %v", url)
				//		message[i] = linebot.NewImageMessage(url, url)
				//	}
				//	if _, err = bot.ReplyMessage(event.ReplyToken, message).Do(); err != nil {
				//		log.Print(err)
				//	}
				} else if message.Text == "長度" {
					if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("beauty len :"+strconv.Itoa(crawler.ImageLengh)+"decard len :"+strconv.Itoa(crawler.DcardImageLengh))).Do(); err != nil {
						log.Print(err)
					}
				}

			}
		}
	}
}
