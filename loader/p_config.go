package loader

import (
	"github.com/IOPaper/Paper/config"
	"github.com/RealFax/pkg-ctl"
)

func init() {
	pkgCtl.RegisterHandler(10, config.LoaderName, config.NewLoader)
}
