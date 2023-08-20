package main

import (
	config2 "github.com/dipdup-io/celestia-indexer/pkg/indexer/config"
	"github.com/dipdup-net/go-lib/config"
)

type Config struct {
	config.Config `yaml:",inline"`
	LogLevel      string         `validate:"omitempty,oneof=debug trace info warn error fatal panic" yaml:"log_level"`
	Indexer       config2.Config `yaml:"indexer"`
}

// Substitute -
func (c *Config) Substitute() error {
	if err := c.Config.Substitute(); err != nil {
		return err
	}
	return nil
}
