package paper

import (
	"context"
	"github.com/IOPaper/Paper/config"
	"github.com/IOPaper/Paper/paper/engine"
	"github.com/IOPaper/Paper/paper/engine/fs"
	"github.com/pkg/errors"
)

func New(ctx context.Context) (engine.Engine, error) {
	assert, err := config.Assert(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "paper.paper")
	}
	switch assert.PaperEngine.Engine {
	case config.PaperEngineTypeFs:
		return fs.NewEngine(ctx, assert.PaperEngine.Fs.Repo)
	case config.PaperEngineTypeSQLite:
		panic("not implement")
	default:
		panic("unsupported engine type")
	}
}
