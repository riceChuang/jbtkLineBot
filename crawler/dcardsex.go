package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/riceChuang/jbtkLineBot/boltdb"
	"github.com/riceChuang/jbtkLineBot/config"
	"github.com/riceChuang/jbtkLineBot/model"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	httpUrl "net/url"
	"strconv"
	"strings"
	"sync"
)

var (
	dcardCraw        *DcardCrawler
	dcardCrawlerOnce = &sync.Once{}
)

type DcardCrawler struct {
	db             *boltdb.Boltdb
	imageLength    int32
	maxImageLength int32
}

func NewDcardCrawler(db *boltdb.Boltdb) *DcardCrawler {
	if dcardCraw != nil {
		return dcardCraw
	}
	dcardCrawlerOnce.Do(func() {
		cfg := config.GetConfig()
		dcardCraw = &DcardCrawler{
			db:             db,
			maxImageLength: cfg.MaxDcardLen,
		}
	})
	return dcardCraw
}

func (d *DcardCrawler) RunCrawlerImage(url string) {
	go d.GetDcarUrl(url)
}

func (d *DcardCrawler) GetImageLength() int32 {
	return d.imageLength
}

func (d *DcardCrawler) GetDcarUrl(url string) {

	client := &http.Client {}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		logrus.Errorf("Dcard NewRequest Error:%v", err)
		return
	}

	res, err := client.Do(req)
	if err != nil {
		logrus.Errorf("Dcard Do Error:%v", err)
		return
	}
	defer res.Body.Close()
	logrus.Println("resp:%+v",res)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logrus.Errorf("Dcard ReadAll Error:%v", err)
		return
	}

	fmt.Println(string(body))

	result := []*model.Dcard{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		logrus.Errorf("Dcard parse Error:%v", err)
	}

	for i, value := range result {

		if d.imageLength > d.maxImageLength {
			return
		}

		if i < 2 && strings.Contains(url, "true") {
			continue
		}

		if value.Gender == "F" && value.LikeCount >= 50 {
			for _, urlValue := range value.Media {
				d.imageLength++
				beautyKey := fmt.Sprintf("dcard-%d", d.imageLength)
				if d.imageLength%50 == 0 {
					logrus.Println("dcard len : %d", d.imageLength)
				}

				if !strings.Contains(urlValue.Url, "https") {
					urlValue.Url = strings.Replace(urlValue.Url, "http", "https", -1)
				}

				dcardInfo, err := json.Marshal(model.DcardInfo{
					ID:    value.ID,
					Image: urlValue.Url,
					Link:  fmt.Sprintf("https://www.dcard.tw/f/sex/p/%v-%v", value.ID, httpUrl.QueryEscape(value.Title)),
				})
				if err != nil {
					logrus.Errorf("dcard error:%v", err)
				}
				d.db.Insert(beautyKey, string(dcardInfo))
			}
		}

		if i+1 == len(result) {
			if strings.Contains(url, "true") {
				newUrl := strings.Replace(url, "true", "false", -1)
				d.GetDcarUrl(newUrl)
			} else if strings.Contains(url, "false") && strings.Contains(url, "before") {
				beforeLocation := strings.Index(url, "before")
				newUrl := url[0 : beforeLocation+7]
				newUrl = fmt.Sprintf("%v%v", newUrl, strconv.Itoa(value.ID))
				d.GetDcarUrl(newUrl)
			} else if strings.Contains(url, "false") {
				newUrl := fmt.Sprintf("%v%v%v", url, "&before=", strconv.Itoa(value.ID))
				d.GetDcarUrl(newUrl)
			}
		}

	}

}
