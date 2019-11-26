package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/riceChuang/jbtkLineBot/boltdb"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	httpUrl "net/url"
)

type DcardCrawler struct {
	db *boltdb.Boltdb
}

type Dcard struct {
	ID        int              `json:"id"`
	Media     []*dcardImageUrl `json:"media"`
	Gender    string           `json:"gender"`
	LikeCount int              `json:"likeCount"`
	Title     string           `json:"title"`
}

type DcardInfo struct {
	ID    int    `json:"id"`
	Image string `json:"image"`
	Link  string `json:"link"`
	Title string `json:"title"`
}

type dcardImageUrl struct {
	Url string `json:"url"`
}

var (
	DcardImageLengh int
)

func NewDcrdCrawler(db *boltdb.Boltdb) *DcardCrawler {
	b := &DcardCrawler{
		db: db,
	}
	return b
}

func (d *DcardCrawler) RunImage(url string) {
	d.GetDcarUrl(url)
}

func (d *DcardCrawler) GetDcarUrl(url string) {
	resp, err := http.Get(url)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	result := []*Dcard{}
	err = json.Unmarshal(body, &result)
	fmt.Print(string(body))
	if err != nil {
		fmt.Print(err)
	}

	for i, value := range result {

		if DcardImageLengh > 600 {
			return
		}

		if i < 2 && strings.Contains(url, "true") {
			continue
		}

		if value.Gender == "F" && value.LikeCount >= 50 {
			for _, urlValue := range value.Media {
				DcardImageLengh++
				beautyKey := fmt.Sprintf("dcard-%d", DcardImageLengh)
				if DcardImageLengh%50 == 0 {
					fmt.Println("dcard len : %d", DcardImageLengh)
				}

				if !strings.Contains(urlValue.Url, "https") {
					urlValue.Url = strings.Replace(urlValue.Url, "http", "https", -1)
				}

				dcardInfo, err := json.Marshal(DcardInfo{
					ID:    value.ID,
					Image: urlValue.Url,
					Link:  fmt.Sprintf("https://www.dcard.tw/f/sex/p/%v-%v", value.ID, httpUrl.QueryEscape(value.Title)),
				})
				if err != nil {
					fmt.Println(err)
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
