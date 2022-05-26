package mock

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"sync"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lesismal/nbio/extension/tls"
	"github.com/lesismal/nbio/logging"
	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"github.com/lesismal/nbio/taskpool"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"eim/global"
	"eim/internal/protocol"
	"eim/model"
	"eim/proto/pb"
)

var connectedCount = &atomic.Int64{}
var sentCount = &atomic.Int64{}
var ackCount = &atomic.Int64{}
var msgCount = &atomic.Int64{}
var hbCount = &atomic.Int64{}
var conns sync.Map
var connectedConsume time.Duration

func Do() {
	var connectionStart, sendStart time.Time
	logging.SetLevel(logging.LevelNone)

	go func() {
		var ticker = time.NewTicker(time.Second)
		var so, ro, sl, rl int64
		lastSent := &atomic.Int64{}
		lastRecv := &atomic.Int64{}

		for {
			select {
			case <-ticker.C:
				{
					sc := sentCount.Load()
					if sc > 0 && lastSent.Load() < sc {
						sl++
						so = sc / sl
					}
					mc := msgCount.Load()
					if mc > 0 && lastRecv.Load() < mc {
						rl++
						ro = mc / rl
					}

					lastSent.Store(sc)
					lastRecv.Store(mc)

					t := table.NewWriter()
					t.SetOutputMirror(os.Stdout)
					t.AppendHeader(table.Row{"Devices", "Send", "Send Tps", "Received", "Received Tps", "Ack", "Heartbeat", "Goroutines"})
					t.AppendRows([]table.Row{{
						connectedCount.Load(), sentCount.Load(), so, msgCount.Load(), ro, ackCount.Load(), hbCount.Load(), runtime.NumGoroutine(),
					}})
					t.AppendSeparator()
					t.Render()

					if sentCount.Load() == msgCount.Load() && msgCount.Load() != 0 {
						global.Logger.Info(fmt.Sprintf("Mock completed，Connection %v : %v，Send %v : %v",
							global.SystemConfig.Mock.ClientCount,
							connectedConsume,
							global.SystemConfig.Mock.ClientCount*global.SystemConfig.Mock.MessageCount,
							time.Since(sendStart)))
						return
					}
				}
			}
		}
	}()

	engine := nbhttp.NewEngine(nbhttp.Config{})
	err := engine.Start()
	if err != nil {
		fmt.Printf("nbio.Start failed: %v\n", err)
		return
	}

	wg := sync.WaitGroup{}
	pool := taskpool.New(runtime.NumCPU(), time.Minute)

	connectionStart = time.Now()

	for i := 0; i < global.SystemConfig.Mock.ClientCount; i++ {
		wg.Add(1)
		pool.Go(func(i int) func() {
			return func() {
				defer wg.Done()
				id := uuid.New().String()
				u := url.URL{Scheme: "ws", Host: global.SystemConfig.Mock.EimEndpoints.Value()[i%len(global.SystemConfig.Mock.EimEndpoints.Value())], Path: "/"}
				var userId, deviceId, deviceName, deviceVersion, deviceType string

				userId = "user_" + id
				deviceId = "device_" + id
				deviceName = "linux_" + id
				deviceVersion = "1.0.0"
				deviceType = model.LinuxDevice

				dialer := &websocket.Dialer{
					Engine:   engine,
					Upgrader: newUpgrader(),
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				}

				var conn *websocket.Conn
				for {
					conn, _, err = dialer.Dial(u.String(), http.Header{
						"UserId":        []string{userId},
						"DeviceId":      []string{deviceId},
						"DeviceName":    []string{deviceName},
						"DeviceVersion": []string{deviceVersion},
						"DeviceType":    []string{deviceType},
					})
					if err != nil {
						global.Logger.Error("Error connecting remote server", zap.String("endpoint", u.String()), zap.Error(err))
						time.Sleep(time.Second)
						continue
					}
					break
				}

				cli := &client{
					userId:        userId,
					deviceId:      deviceId,
					deviceName:    deviceName,
					deviceType:    deviceType,
					deviceVersion: deviceVersion,
					conn:          conn,
				}

				conn.SetSession(cli)
				conns.Store(cli.userId, cli)

				connectedCount.Add(1)
			}
		}(i))
	}

	wg.Wait()

	connectedConsume = time.Since(connectionStart)

	sendStart = time.Now()

	conns.Range(func(key, value interface{}) bool {
		go func(cli *client) {
			var msgTotal = global.SystemConfig.Mock.MessageCount
			ticker := time.NewTicker(time.Second)
			for {
				select {
				case <-ticker.C:
					{
						msg := &pb.Message{
							MsgId:      uuid.New().String(),
							MsgType:    model.TextMessage,
							Content:    time.Now().String(),
							FromType:   model.FromUser,
							FromId:     cli.userId,
							FromName:   cli.userId,
							FromDevice: cli.deviceId,
							ToType:     model.ToUser,
							ToId:       cli.userId,
							ToName:     cli.userId,
							ToDevice:   cli.deviceId,
						}
						body, _ := proto.Marshal(msg)
						err := cli.conn.WriteMessage(websocket.BinaryMessage, protocol.WebsocketCodec.Encode(protocol.Message, body))
						if err != nil {
							global.Logger.Warn("Error sending message", zap.Error(err))
							return
						}
						sentCount.Add(1)
						msgTotal--
						if msgTotal == 0 {
							return
						}
					}
				}
			}
		}(value.(*client))

		go func(cli *client) {
			ticker := time.NewTicker(time.Second * 30)
			for {
				select {
				case <-ticker.C:
					{
						_ = cli.conn.WriteMessage(websocket.PingMessage, []byte(time.Now().String()))
					}
				}
			}
		}(value.(*client))
		return true
	})
}
