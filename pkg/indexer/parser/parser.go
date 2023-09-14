package parser

import (
	"context"
	"github.com/dipdup-net/indexer-sdk/pkg/modules"
)

type Module struct {
	modules.BaseModule
}

var _ modules.Module = (*Module)(nil)

const (
	InputName  = "blocks"
	OutputName = "data"
	StopOutput = "stop"
)

func NewModule() Module {
	return Module{
		BaseModule: modules.New("parser"),
	}
}

func (p *Module) Start(ctx context.Context) {
	p.Log.Info().Msg("starting parser module...")
	p.G.GoCtx(ctx, p.listen)
}

func (p *Module) Close() error {
	p.Log.Info().Msg("closing...")
	p.G.Wait()
	return nil
}
