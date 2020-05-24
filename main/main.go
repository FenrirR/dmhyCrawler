package main

import "dmhyCrawler/crawler"

func main() {
	myCrawler := crawler.NewCrawler("某科学超电磁炮T",
		"https://share.dmhy.org/topics/list?keyword=lolihouse+%E7%94%B5%E7%A3%81%E7%82%AE")
	myCrawler.FetchMagLinks()
}
