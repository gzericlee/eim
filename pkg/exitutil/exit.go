package exitutil

import (
	"os"
	"os/signal"
)

func WaitSignal(callback func()) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	callback()
}
