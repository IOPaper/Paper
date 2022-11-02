package paper

import (
	"context"
	"github.com/IOPaper/Paper/config"
	"github.com/IOPaper/Paper/paper/engine"
	"github.com/IOPaper/Paper/paper/engine/fs"
)

func New(ctx context.Context, conf *config.Config) (engine.Engine, error) {
	switch conf.PaperEngine.Engine {
	case config.PaperEngineTypeFs:
		return fs.NewEngine(ctx, conf.PaperEngine.Fs.Repo)
	case config.PaperEngineTypeSQLite:
		panic("not implement")
	default:
		panic("unsupported engine type")
	}
}
