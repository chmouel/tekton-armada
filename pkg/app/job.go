package app

type Job struct {
	Tags    []string          `json:"tags"`
	Volumes map[string]string `json:"volumes"`
	Yamls   []string          `json:"yamls"`
}
