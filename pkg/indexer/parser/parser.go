package parser

import (
	"context"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"sync"
)

type Parser struct {
	input  *modules.Input
	output *modules.Output
	log    zerolog.Logger
	wg     *sync.WaitGroup
}

const (
	name        = "parser"
	BlocksInput = "blocks"
	DataOutput  = "data"
)

func NewModule() Parser {
	return Parser{
		input:  modules.NewInput(BlocksInput),
		output: modules.NewOutput(DataOutput),
		log:    log.With().Str("module", name).Logger(),
		wg:     new(sync.WaitGroup),
	}
}

// Name -
func (*Parser) Name() string {
	return name
}

func (p *Parser) Start(ctx context.Context) {
	p.log.Info().Msg("starting parser module...")

	p.wg.Add(1)
	go p.listen(ctx)
}
