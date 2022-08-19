package boot

import (
	"github.com/IOPaper/Paper/boot/ctl"
	"github.com/IOPaper/Paper/paper"

	"github.com/IOPaper/Paper/global"
	"github.com/IOPaper/Paper/service/web"
)

var registry = []ctl.Created{
	{"config-loader", global.NewConfig},
	{"paper-engine-loader", paper.New},
	{"web-server", web.New},
}

func register() (err error) {
	for _, created := range registry {
		if err = ctl.Register(created.Name, created.Func); err != nil {
			return err
		}
	}
	return nil
}
