package crawler

import "github.com/riceChuang/jbtkLineBot/boltdb"

var(
	beautyCrawler *BeautyCrawler
	dcardSex *DcardCrawler
	joker *JokerCrawler
)

func Initialize(){
	beautyCrawler = NewBeautyCrawler(boltdb.DB())
	dcardSex = NewDcrdCrawler(boltdb.DB())
	joker = NewJokerCrawler(boltdb.DB())
}

