package main

import (
	"dmhyCrawler/conf"
	"dmhyCrawler/crawler"
)

func main() {
	serviceConf := conf.GetServiceConfig()
	for _, bangumiConfig := range serviceConf.BangumiConfigs {
		url := "https://share.dmhy.org/topics/list?keyword=" + bangumiConfig.SearchExp
		myCrawler := crawler.NewCrawler(bangumiConfig.Title, url)
		myCrawler.FetchMagLinks()

	}
}
