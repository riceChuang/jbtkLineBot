package service

import (
	"encoding/json"
	"fmt"
	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/riceChuang/jbtkLineBot/boltdb"
	"github.com/riceChuang/jbtkLineBot/crawler"
	"github.com/riceChuang/jbtkLineBot/line"
	"github.com/riceChuang/jbtkLineBot/model"
	"github.com/sirupsen/logrus"
	"log"
	"math/rand"
)

func ReplyDcardMapMessage(event *linebot.Event) {
	imageColumns := []*linebot.ImageCarouselColumn{}
	lineClient := line.GetBotClient()
	instance, err := crawler.GetCrawlerByType(crawler.DcardSex)
	if err != nil {
		logrus.Error("cant find instance")
		return
	}
	for i := 0; i < 10; i++ {
		imageIndex := rand.Intn(int(instance.GetImageLength()))
		db := boltdb.DB()
		dbKey := fmt.Sprintf("dcard-%d", imageIndex)
		dcardResp := db.Read(dbKey)
		dcardInfo := model.DcardInfo{}
		err := json.Unmarshal([]byte(dcardResp), &dcardInfo)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(dcardInfo)
		imageColumns = append(imageColumns, linebot.NewImageCarouselColumn(dcardInfo.Image, linebot.NewURIAction(dcardInfo.Title, dcardInfo.Link)))
	}

	if _, err := lineClient.ReplyMessage(event.ReplyToken, linebot.NewTemplateMessage("Dacrd", linebot.NewImageCarouselTemplate(imageColumns...))).Do(); err != nil {
		log.Print(err)
	}
}

func ReplyDcardMessage(event *linebot.Event) {
	instance, err := crawler.GetCrawlerByType(crawler.DcardSex)
	if err != nil {
		logrus.Error("cant find instance")
		return
	}
	imageIndex := rand.Intn(int(instance.GetImageLength()))
	lineClient := line.GetBotClient()
	db := boltdb.DB()
	dbkey := fmt.Sprintf("dcard-%d", imageIndex)
	dcardResp := db.Read(dbkey)
	dcardInfo := model.DcardInfo{}
	err = json.Unmarshal([]byte(dcardResp), &dcardInfo)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("my image link: %v", dcardInfo.Image)
	if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(dcardInfo.Image, dcardInfo.Image)).Do(); err != nil {
		log.Print(err)
	}
}

func ReplyPttMessage(event *linebot.Event) {
	instance, err := crawler.GetCrawlerByType(crawler.Beauty)
	if err != nil {
		logrus.Error("cant find instance")
		return
	}
	imageIndex := rand.Intn(int(instance.GetImageLength()))
	lineClient := line.GetBotClient()
	db := boltdb.DB()
	dbkey := fmt.Sprintf("beauty-%d", imageIndex)
	url := db.Read(dbkey)
	fmt.Printf("my image link: %v", url)
	if _, err := lineClient.ReplyMessage(event.ReplyToken, linebot.NewImageMessage(url, url)).Do(); err != nil {
		log.Print(err)
	}
}

func ReplyJokeMessage(event *linebot.Event) {
	instance, err := crawler.GetCrawlerByType(crawler.Joker)
	if err != nil {
		logrus.Error("cant find instence")
		return
	}
	imageIndex := rand.Intn(int(instance.GetImageLength()))
	lineClient := line.GetBotClient()
	dbkey := fmt.Sprintf("joker-%d", imageIndex)
	content := crawler.JokerMap[dbkey]
	if _, err := lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(content)).Do(); err != nil {
		log.Print(err)
	}
}

func ReplyMessageLen(event *linebot.Event) {
	lineClient := line.GetBotClient()
	jokerInstance, err := crawler.GetCrawlerByType(crawler.Joker)
	if err != nil {
		logrus.Error("cant find jokerInstance")
		return
	}
	dcardInstance, err := crawler.GetCrawlerByType(crawler.DcardSex)
	if err != nil {
		logrus.Error("cant find dcardInstance")
		return
	}
	beautyInstance, err := crawler.GetCrawlerByType(crawler.Beauty)
	if err != nil {
		logrus.Error("cant find beautyInstance")
		return
	}

	respMessage := fmt.Sprintf("beautyLen: %v, dcardLen: %v, jokerLen:%v", beautyInstance.GetImageLength(), dcardInstance.GetImageLength(), jokerInstance.GetImageLength())
	if _, err = lineClient.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(respMessage)).Do(); err != nil {
		log.Print(err)
	}
}

func ReplyPornMessage(event *linebot.Event) {

}
