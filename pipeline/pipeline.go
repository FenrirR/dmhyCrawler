package pipeline

import (
	"dmhyCrawler/utils"
	"sync"
)

var resDir = "./resource/links/"
var histDir = "./resource/history/"

type IPipeline interface {
	process(item Item) (abort bool)
}

type persistPipeline struct {
	sync.RWMutex
}

func initPersistPipeline() *persistPipeline {
	utils.CreateDir(resDir)
	utils.RemoveContents(resDir)
	return &persistPipeline{}
}

func (p *persistPipeline) process(item Item) (abort bool) {
	p.Lock()
	defer p.Unlock()
	path := resDir + item.Title + ".txt"
	utils.CreateFile(path)
	utils.SaveList2Txt([]string{item.MagLink}, path)
	return
}

type dedupPipeline struct {
	sync.RWMutex
}

func initDedupPipeline() *dedupPipeline {
	utils.CreateDir(histDir)
	return &dedupPipeline{}
}

func (p *dedupPipeline) process(item Item) (abort bool) {
	p.Lock()
	defer p.Unlock()
	path := histDir + item.Title + ".txt"
	utils.CreateFile(path)
	histData := utils.ReadTxt2Set(path)

	if histData[item.MagLink] {
		abort = true
	} else {
		utils.SaveList2Txt([]string{item.MagLink}, path)
	}
	return
}
