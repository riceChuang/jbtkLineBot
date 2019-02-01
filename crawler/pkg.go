package crawler

import "github.com/riceChuang/jbtkLineBot/boltdb"

var(
	_b *BeautyCrawler
)

func Initialize(){
	_b = NewBeautyCrawler(boltdb.DB())
}

