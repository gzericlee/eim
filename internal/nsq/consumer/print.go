package consumer

import (
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/nsqio/go-nsq"
)

func printConsumersDetail() {
	go func() {
		ticker := time.NewTicker(time.Second * 5)
		for {
			select {
			case <-ticker.C:
				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"Consumer", "Connections", "Received", "Finished", "Requeued"})
				consumers.Range(func(key, value interface{}) bool {
					stats := value.(*nsq.Consumer).Stats()
					t.AppendRows([]table.Row{{
						key,
						stats.Connections,
						stats.MessagesReceived,
						stats.MessagesFinished,
						stats.MessagesRequeued},
					})
					return true
				})
				t.AppendSeparator()
				t.Render()
			}
		}
	}()
}
