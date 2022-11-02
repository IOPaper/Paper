package fs

import (
	"context"
	"github.com/IOPaper/Paper/paper/engine"
)

type Engine struct {
	fapi FileApi
	ctx  context.Context
	root string
}

func NewEngine(ctx context.Context, root string) engine.Engine {
	return &Engine{
		fapi: NewFileApi(root),
		ctx:  ctx,
		root: root,
	}
}

func (e *Engine) Type() string {
	return "fs"
}

func (e *Engine) GetOnePaper(key string) (*engine.Paper, error) {
	return nil, nil
}
