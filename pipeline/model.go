package pipeline

// Step представляет собой шаг пайплайна
type Step struct {
	Name     string   `yaml:"name"`
	Image    string   `yaml:"image"`
	Commands []string `yaml:"commands"`
}

// Pipeline представляет собой структуру пайплайна CI/CD
type Pipeline struct {
	Steps []Step `yaml:"steps"`
}

// PipelineRawString структура пайплайна в string
type PipelineRawString struct {
	Content string `json:"content"`
}
