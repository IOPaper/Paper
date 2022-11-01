package main

import (
	"context"
	"github.com/IOPaper/Paper/loader"
	"github.com/RealFax/pkg-ctl"
)

func main() {

	pkgCtl.SetupActive(loader.Loader)

	var (
		ctx, cancel = context.WithCancel(context.Background())
		err         error
	)

	if err = pkgCtl.Startup(&ctx); err != nil {
		return
	}

	if err = pkgCtl.ListenAndDestroy(cancel); err != nil {
		return
	}

}
