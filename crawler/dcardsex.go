package crawler

import (
	"encoding/json"
	"fmt"
	"github.com/riceChuang/jbtkLineBot/boltdb"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type DcardCrawler struct {
	db *boltdb.Boltdb
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

	result := []*dcard{}
	err = json.Unmarshal(body, &result)
	fmt.Print(string(body))
	if err != nil {
		fmt.Print(err)
	}

	for i, value := range result {

		if DcardImageLengh > 1000 {
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

				d.db.Insert(beautyKey, urlValue.Url)
			}
		}

		if i+1 == len(result) {
			if strings.Contains(url, "true") {
				newUrl := strings.Replace(url, "true", "false", -1)
				d.GetDcarUrl(newUrl)
			} else if strings.Contains(url, "false") && strings.Contains(url, "before") {
				beforeLocation := strings.Index(url, "before")
				newUrl := url[0 : beforeLocation+7]
				newUrl = fmt.Sprintf("%v%v", newUrl, strconv.Itoa(value.Id))
				d.GetDcarUrl(newUrl)
			} else if strings.Contains(url, "false") {
				newUrl := fmt.Sprintf("%v%v%v", url, "&before=", strconv.Itoa(value.Id))
				d.GetDcarUrl(newUrl)
			}
		}

	}

}

type dcard struct {
	Id        int              "json:id"
	Media     []*dcardImageUrl "json:media"
	Gender    string           "json:gender"
	LikeCount int              "json:likeCount"
}

type dcardImageUrl struct {
	Url string "json:url"
}
