package service

import (
	"encoding/json"
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/riceChuang/jbtkLineBot/boltdb"
	"github.com/riceChuang/jbtkLineBot/crawler"
	"log"
	"math/rand"
)

const (
	MessagePtt    = "抽"
	MessageJoke   = "笑"
	MessageDcard  = "d"
	MessageDcards = "dd"
	MessageLong   = "長度"
	MessagePorn   = "片"
	MessageTest   = "測"
)

func ReplyDcardMapMessage(event *linebot.Event) {
	imageColumns := []*linebot.ImageCarouselColumn{}

	for i := 0; i < 10; i++ {
		imageIndex := rand.Intn(crawler.DcardImageLengh)
		db := boltdb.DB()
		dbkey := fmt.Sprintf("dcard-%d", imageIndex)
		dcardResp := db.Read(dbkey)
		dcardInfo := crawler.DcardInfo{}
		err := json.Unmarshal([]byte(dcardResp), &dcardInfo)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(dcardInfo)
		imageColumns = append(imageColumns, linebot.NewImageCarouselColumn(dcardInfo.Image, linebot.NewURIAction(dcardInfo.Title, dcardInfo.Link)))
	}
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("Dacrd", linebot.NewImageCarouselTemplate(imageColumns...))).Do(); err != nil {
		log.Print(err)
	}
}

func ReplyDcardMessage(event *linebot.Event) {
	imageIndex := rand.Intn(crawler.DcardImageLengh)
	db := boltdb.DB()
	dbkey := fmt.Sprintf("dcard-%d", imageIndex)
	dcardResp := db.Read(dbkey)
	dcardInfo := crawler.DcardInfo{}
	err := json.Unmarshal([]byte(dcardResp), &dcardInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("my image link: %v", dcardInfo.Image)
	if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(dcardInfo.Image, dcardInfo.Image)).Do(); err != nil {
		log.Print(err)
	}
}

func ReplyPttMessage(event *linebot.Event) {
	imageIndex := rand.Intn(crawler.ImageLength)
	db := boltdb.DB()
	dbkey := fmt.Sprintf("beauty-%d", imageIndex)
	url := db.Read(dbkey)
	fmt.Printf("my image link: %v", url)
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(url, url)).Do(); err != nil {
		log.Print(err)
	}
}

func ReplyJokeMessage(event *linebot.Event) {
	imageIndex := rand.Intn(crawler.JokerLenght)
	dbkey := fmt.Sprintf("joker-%d", imageIndex)
	content := crawler.JokerMap[dbkey]
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(content)).Do(); err != nil {
		log.Print(err)
	}
}

func ReplyPornMessage(event *linebot.Event) {

}

func ReplyTransferImage(event *linebot.Event) {
	imageURL := "https://i.imgur.com/Cj28dWe.jpg"
	if _, err := bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("yoyoyo").WithQuickReplies(linebot.NewQuickReplyItems(linebot.NewQuickReplyButton(imageURL, linebot.NewCameraRollAction("上傳嗎"))))).Do(); err != nil {
		log.Print(err)
	}
}
