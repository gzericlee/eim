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
				t.AppendHeader(table.Row{"devices", "L7 CPS", "Received", "Send", "Ack", "Invalid", "Error", "Heartbeat", "Goroutines"})
				t.AppendRows([]table.Row{{
					its.clientTotal,
					maxL7Cps,
					its.receivedMsgTotal,
					its.sendMsgTotal,
					its.ackTotal,
					its.invalidMsgTotal,
					its.errorTotal,
					its.heartbeatTotal,
					runtime.NumGoroutine()},
				})
				t.AppendSeparator()
				t.Render()

				mMetric, _ := metric.GetMachineMetric()
				err := its.storageRpc.RegisterGateway(&model.Gateway{
					Ip:             config.SystemConfig.LocalIp,
					ClientTotal:    its.clientTotal,
					SendTotal:      its.sendMsgTotal,
					ReceivedTotal:  its.receivedMsgTotal,
					InvalidTotal:   its.invalidMsgTotal,
					GoroutineTotal: int64(runtime.NumGoroutine()),
					MemUsed:        float32(mMetric.MemUsed),
					CpuUsed:        float32(mMetric.CpuUsed),
				}, time.Second*10)
				if err != nil {
					return
				}
			}
		}
	}
}
