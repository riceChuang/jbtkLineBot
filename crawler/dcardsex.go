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
	//
	//client := &http.Client {}
	//req, err := http.NewRequest(http.MethodGet, url, nil)
	//if err != nil {
	//	logrus.Errorf("Dcard NewRequest Error:%v", err)
	//	return
	//}
	//
	//res, err := client.Do(req)
	//if err != nil {
	//	logrus.Errorf("Dcard Do Error:%v", err)
	//	return
	//}
	//defer res.Body.Close()
	//logrus.Println("resp:%+v",res)
	//body, err := ioutil.ReadAll(res.Body)
	//if err != nil {
	//	logrus.Errorf("Dcard ReadAll Error:%v", err)
	//	return
	//}


	method := "GET"

	client := &http.Client {
	}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("authority", " www.dcard.tw'")
	req.Header.Add("pragma", " no-cache'")
	req.Header.Add("cache-control", " no-cache'")
	req.Header.Add("sec-ch-ua", " \" Not;A Brand\";v=\"99\", \"Google Chrome\";v=\"91\", \"Chromium\";v=\"91\"'")
	req.Header.Add("sec-ch-ua-mobile", " ?0'")
	req.Header.Add("user-agent", " Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.114 Safari/537.36'")
	req.Header.Add("accept", " */*'")
	req.Header.Add("sec-fetch-site", " same-origin'")
	req.Header.Add("sec-fetch-mode", " cors'")
	req.Header.Add("sec-fetch-dest", " empty'")
	req.Header.Add("referer", " https://www.dcard.tw/f'")
	req.Header.Add("accept-language", " zh-TW,zh;q=0.9,en-US;q=0.8,en;q=0.7'")
	req.Header.Add("Cookie", "__gads=ID=93420f7732023247:T=1615971361:S=ALNI_MYiV4izBJAvHP3lyZjddGz_SzZKZQ; __auc=8bad7fe017911f8edb8b9b83600; _gid=GA1.2.1941835643.1626374400; dcsrd=YH-KXakKJ3HqDvpCjSlp2qi9; CFFPCKUUID=8602-ZjO4cpjNHiQpK7lgXvSULfJJ0s1cRWr7; CFFPCKUUIDMAIN=8157-oHkzb6QhaslU5pDIp5mE0pbvFeLbUA9F; __htid=f9b901d9-216a-4435-b517-d8db3579ad40; _ht_50ef57=1; __asc=cfba216117aad6ac46a11d0802e; _gat=1; _ga_C3J49QFLW7=GS1.1.1626407095.14.1.1626408512.0; _ga=GA1.2.1289658169.1615971361; dcard-web-oauth-cv=T3RtZ0swR3V6ZDNBUnI0bGs0QWxUcnkxTkFBeE0wcWJMMlJ0a0xTbjF0TQ==; dcard-web-oauth-cv.sig=mLg3Z590Hz3vI3XPiHsx9uvEyXc; __cf_bm=356262476654de78126e9860e81949bfb7b1b0cd-1626408514-1800-ATtHS8GoqlmVCgI2xfOQIVbj5T/NR/vs6nh6PAMcRQhnyYnPOpyr0Nu1iIZrC1T/W/LILbTO611XSpe3YKbBMBo=")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))

	logrus.Printf("this is Body :%s",string(body))
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
