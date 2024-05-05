package pipeline

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

func parsePipeline(data []byte) (Pipeline, error) {
	pipeline := Pipeline{}

	err := yaml.Unmarshal(data, pipeline)
	if err != nil {
		log.Error("parsePipeline #0: ", err)
		return pipeline, err
	}

	return pipeline, nil
}
