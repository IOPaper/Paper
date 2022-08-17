package boot

import (
	"github.com/IOPaper/Paper/boot/ctl"
	"os"
	"os/signal"
	"syscall"
)

func Start() (err error) {
	if err = register(); err != nil {
		return
	}
	return ctl.Startup()
}

// Shutdown listen to and intercept the system exit signal, and then call the destroy func
func Shutdown() error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	<-c
	return ctl.Destroy()
}
