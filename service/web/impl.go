package web

import (
	"context"
	"github.com/IOPaper/Paper/boot/ctl"
	"github.com/IOPaper/Paper/global"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type Implement struct {
	route  *gin.Engine
	server *http.Server
}

func New() ctl.I {
	return &Implement{}
}

func (i *Implement) Create() error {

	gin.SetMode(global.Config.Engine.LogLevel.String())

	i.route = gin.Default()

	_ = i.route.SetTrustedProxies(nil)

	return nil
}

func (i *Implement) Destroy() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	return i.server.Shutdown(ctx)
}

func (i *Implement) Start() error {
	i.server = &http.Server{
		Addr:    global.Config.Http.Addr,
		Handler: i.route,
	}
	return i.server.ListenAndServe()
}

func (i *Implement) IsAsync() bool {
	return true
}
