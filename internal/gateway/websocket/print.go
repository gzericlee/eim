package websocket

import (
	"os"
	"runtime"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	"eim/internal/config"
	"eim/internal/metric"
	"eim/internal/model"
)

var maxL7Cps, lastClientTotal int64

func (its *Server) printServiceDetail() {
	ticker := time.NewTicker(time.Second * 5)
	for {
		select {
		case <-ticker.C:
			{
				if its.clientTotal > lastClientTotal {
					l7Cps := its.clientTotal - lastClientTotal
					if l7Cps > maxL7Cps {
						maxL7Cps = l7Cps
					}
					lastClientTotal = its.clientTotal
				}
				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"Devices", "L7 CPS", "Received", "Send", "Invalid", "Error", "Heartbeat", "Goroutines"})
				t.AppendRows([]table.Row{{
					its.clientTotal,
					maxL7Cps,
					its.receivedMsgTotal,
					its.sendMsgTotal,
					its.invalidMsgTotal,
					its.errorTotal,
					its.heartbeatTotal,
					runtime.NumGoroutine()},
				})
				t.AppendSeparator()
				t.Render()

				mMetric, _ := metric.GetMachineMetric()
				err := its.redisManager.RegisterGateway(&model.Gateway{
					Ip:               config.SystemConfig.LocalIp,
					ClientTotal:      its.clientTotal,
					SendMsgTotal:     its.sendMsgTotal,
					ReceivedMsgTotal: its.receivedMsgTotal,
					InvalidMsgTotal:  its.invalidMsgTotal,
					GoroutineTotal:   int64(runtime.NumGoroutine()),
					MemUsed:          float32(mMetric.MemUsed),
					CpuUsed:          float32(mMetric.CpuUsed),
				})
				if err != nil {
					return
				}
			}
		}
	}
}
