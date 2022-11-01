package config

import (
	"context"
	"github.com/RealFax/pkg-ctl"
)

const LoaderName string = "config-loader"

type Implement struct {
	ctx *context.Context
}

func NewLoader(ctx *context.Context) pkgCtl.Handler {
	return &Implement{ctx: ctx}
}

func (i *Implement) Create() error {
	parseArgs(i.ctx)
	return nil
}

func (i *Implement) Start() error {
	return WithYAMLConfig(i.ctx)
}

func (i *Implement) Destroy() error {
	return nil
}

func (i *Implement) IsAsync() bool {
	return false
}
