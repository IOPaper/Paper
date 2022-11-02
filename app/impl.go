package app

import (
	"context"
	"github.com/IOPaper/Paper/app/impls"
	"github.com/IOPaper/Paper/config"
	"github.com/IOPaper/Paper/paper"
	pkgCtl "github.com/RealFax/pkg-ctl"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type Implement struct {
	ctx   context.Context
	route *gin.Engine
	serve *http.Server
}

func NewService(ctx *context.Context) pkgCtl.Handler {
	return &Implement{ctx: *ctx}
}

func (i *Implement) SetupRoute() error {
	engine, err := paper.New(i.ctx)
	if err != nil {
		return err
	}

	rPaper := i.route.Group("/x")
	{

		implPaper := impls.NewImplPaper(engine)
		rPaper.GET("/list", implPaper.GetList)

		rDeep := rPaper.Group("/:index")
		{
			rDeep.GET("/:attachment", implPaper.GetAttachment)
			rDeep.GET("/", implPaper.GetPaper)
		}
	}
	return nil
}

func (i *Implement) Create() error {

	conf, err := config.Assert(i.ctx)
	if err != nil {
		return errors.Wrap(err, "app.impl")
	}

	i.serve = &http.Server{
		Addr:    conf.Http.Addr,
		Handler: i.route,
	}

	gin.SetMode(conf.Http.LogLevel.String())

	i.route = gin.Default()

	_ = i.route.SetTrustedProxies(nil)

	return i.SetupRoute()
}

func (i *Implement) Start() error {
	return i.serve.ListenAndServe()
}

func (i *Implement) Destroy() error {
	return i.serve.Shutdown(context.Background())
}

func (i *Implement) IsAsync() bool {
	return true
}
