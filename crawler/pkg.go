package crawler

import "github.com/riceChuang/jbtkLineBot/boltdb"

var(
	_b *BeautyCrawler
	_dsex *DcardCrawler
)

func Initialize(){
	_b = NewBeautyCrawler(boltdb.DB())
	_dsex = NewDcrdCrawler(boltdb.DB())
}

