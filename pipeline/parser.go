package pipeline

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func parsePipeline(data []byte) (Pipeline, error) {
	pipeline := Pipeline{}

	var pipelineRaw struct{ Pipeline Pipeline }
	err := yaml.Unmarshal(data, &pipelineRaw)
	if err != nil {
		log.Error("parsePipeline #0: ", err)
		return pipeline, err
	}

	pipeline.Steps = pipelineRaw.Pipeline.Steps

	return pipeline, nil
}
