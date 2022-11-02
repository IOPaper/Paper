package loader

import (
	"github.com/IOPaper/Paper/app"
	"github.com/RealFax/pkg-ctl"
)

func init() {
	pkgCtl.RegisterHandler(100, "app-service", app.NewService)
}
