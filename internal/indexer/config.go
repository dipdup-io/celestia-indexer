package indexer

type Config struct {
	Name         string `yaml:"name" validate:"omitempty"`
	Timeout      uint64 `yaml:"timeout" validate:"omitempty"`
	ThreadsCount uint32 `yaml:"threads_count" validate:"omitempty,min=1"`
	Node         *Node  `yaml:"node" validate:"omitempty"`
}

type Node struct {
	Url string `yaml:"url" validate:"omitempty,url"`
	Rps uint64 `yaml:"requests_per_second" validate:"omitempty,min=1"`
}
