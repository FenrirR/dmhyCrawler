package pipeline

var pipelineQueue []IPipeline

func init() {
	configPipeline()
}

func configPipeline() {
	pipelineQueue = append(pipelineQueue, initDedupPipeline())
	pipelineQueue = append(pipelineQueue, initPersistPipeline())
}

func RunPipeline(item Item) {
	for _, pipeline := range pipelineQueue {
		if abort := pipeline.process(item); abort {
			break
		}
	}
}
