package websocket

import (
	"os"
	"runtime"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
)

var maxL7Cps, lastClientTotal int64

func (its *Server) PrintServiceStats() {
	ticker := time.NewTicker(time.Second * 30)
	for {
		select {
		case <-ticker.C:
			{
				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"Devices", "Received", "Send", "Invalid", "Error", "Ack", "Heartbeat", "Goroutines"})
				t.AppendRows([]table.Row{{
					its.clientTotal,
					its.receivedMsgTotal,
					its.sendMsgTotal,
					its.invalidMsgTotal,
					its.errorTotal,
					its.ackTotal,
					its.heartbeatTotal,
					runtime.NumGoroutine()},
				})
				t.AppendSeparator()
				t.Render()
			}
		}
	}
}
