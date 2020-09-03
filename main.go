package main

import (
	"fmt"
	"github.com/riceChuang/jbtkLineBot/boltdb"
	"github.com/riceChuang/jbtkLineBot/crawler"
	"github.com/riceChuang/jbtkLineBot/service"
	"github.com/riceChuang/jbtkLineBot/types"
	"net/http"
	"os"
)

func main() {

	boltdb.Initialize()
	crawler.Initialize()
	types.InitialConfigPkg()
	config := types.GetConfig()

	crawlerTypesMap := map[crawler.Type]string{
		crawler.Beauty:   config.BeautyUrl,
		crawler.DcardSex: config.DcardUrl,
	}

	for crawlerType := range crawlerTypesMap {
		crawlerWorker, err := crawler.GetCrawlerByType(crawlerType)
		if err != nil {
			fmt.Println(err)
		}
		go crawlerWorker.RunImage(crawlerTypesMap[crawlerType])
	}

	service.InitBotClient()
	http.HandleFunc("/back", service.CallbackHandler)
	port := os.Getenv("PORT")
	addr := ""
	if port != "" {
		addr = fmt.Sprintf(":%s", port)
	} else {
		addr = ":3000"
	}

	http.ListenAndServe(addr, nil)
}

//
//func PrintMemUsage() {
//	var m runtime.MemStats
//	runtime.ReadMemStats(&m)
//	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
//	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
//	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
//	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
//	fmt.Printf("\tNumGC = %v\n", m.NumGC)
//}
//func bToMb(b uint64) uint64 {
//	return b / 1024 / 1024
//}
