package websocket

import (
	"os"
	"runtime"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"

	"eim/internal/config"
	"eim/internal/metric"
	"eim/internal/redis"
	"eim/internal/types"
)

var maxL7Cps, lastClientTotal, lastMaxL7Cps int64

func printWebSocketServiceDetail() {
	go func() {
		ticker := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-ticker.C:
				{
					currentClientTotal := gatewaySvr.clientTotal.Load()
					if currentClientTotal > lastClientTotal {
						l7Cps := currentClientTotal - lastClientTotal
						if l7Cps > maxL7Cps {
							maxL7Cps = l7Cps
						}
						lastClientTotal = currentClientTotal
					}
					t := table.NewWriter()
					t.SetOutputMirror(os.Stdout)
					t.AppendHeader(table.Row{"Devices", "L7 CPS", "Received", "Sent", "Invalid", "Heartbeat", "Goroutines"})
					t.AppendRows([]table.Row{{
						gatewaySvr.clientTotal.Load(),
						maxL7Cps,
						gatewaySvr.receivedTotal.Load(),
						gatewaySvr.sentTotal.Load(),
						gatewaySvr.invalidMsgTotal.Load(),
						gatewaySvr.heartbeatTotal.Load(),
						runtime.NumGoroutine()},
					})
					t.AppendSeparator()
					t.Render()

					mMetric, _ := metric.GetMachineMetric()
					err := redis.RegisterGateway(&types.Gateway{
						Ip:             config.SystemConfig.LocalIp,
						ClientTotal:    gatewaySvr.clientTotal.Load(),
						SentTotal:      gatewaySvr.sentTotal.Load(),
						ReceivedTotal:  gatewaySvr.receivedTotal.Load(),
						InvalidTotal:   gatewaySvr.invalidMsgTotal.Load(),
						GoroutineTotal: int64(runtime.NumGoroutine()),
						MemUsed:        mMetric.MemUsed,
						CpuUsed:        mMetric.CpuUsed,
					})
					if err != nil {
						return
					}
				}
			}
		}
	}()
}
