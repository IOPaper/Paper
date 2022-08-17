package web

import (
	"context"
	"github.com/IOPaper/Paper/boot/ctl"
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
		Addr:    "",
		Handler: i.route,
	}
	return nil
}

func (i *Implement) IsAsync() bool {
	return true
}
