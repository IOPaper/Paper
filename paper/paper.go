package paper

import (
	"github.com/IOPaper/Paper/boot/ctl"
	"github.com/IOPaper/Paper/global"
	"github.com/IOPaper/Paper/paper/core"
	"github.com/IOPaper/Paper/paper/engine"
)

var Func core.PaperFunc

type Implement struct {
}

func New() ctl.I {
	return &Implement{}
}

func (i *Implement) Create() error {
	return nil
}

func (i *Implement) Destroy() error {
	return Func.Close()
}

func (i *Implement) Start() (err error) {
	Func, err = engine.New(global.Config.Engine.Repo)
	return
}

func (i *Implement) IsAsync() bool {
	return false
}
