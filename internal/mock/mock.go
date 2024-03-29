package mock

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"runtime"
	"strconv"
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

	"eim/internal/config"
	"eim/internal/protocol"
	"eim/internal/types"
	"eim/pkg/log"
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

					if sentCount.Load() == msgCount.Load() && sentCount.Load() == int64(config.SystemConfig.Mock.ClientCount*config.SystemConfig.Mock.MessageCount) && msgCount.Load() != 0 {
						log.Info(fmt.Sprintf("Mock completed，Connection %v : %v，Send %v : %v",
							config.SystemConfig.Mock.ClientCount,
							connectedConsume,
							config.SystemConfig.Mock.ClientCount*config.SystemConfig.Mock.MessageCount,
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
	pool := taskpool.NewFixedPool(runtime.NumCPU(), 1024)

	connectionStart = time.Now()

	for i := 1; i <= config.SystemConfig.Mock.ClientCount; i++ {
		wg.Add(1)
		pool.Go(func(i int) func() {
			return func() {
				defer wg.Done()
				id := strconv.Itoa(i)
				u := url.URL{Scheme: "ws", Host: config.SystemConfig.Mock.EimEndpoints.Value()[i%len(config.SystemConfig.Mock.EimEndpoints.Value())], Path: "/"}

				token := base64.StdEncoding.EncodeToString([]byte("lirui@bingo:pass@word1"))
				userId := "user-" + id
				deviceId := "device-" + id
				deviceName := "linux-" + id
				deviceVersion := "1.0.0"
				deviceType := types.LinuxDevice

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
						"Token":         []string{token},
						"DeviceId":      []string{deviceId},
						"DeviceName":    []string{deviceName},
						"DeviceVersion": []string{deviceVersion},
						"DeviceType":    []string{deviceType},
					})
					if err != nil {
						log.Error("Error connecting remote server", zap.String("endpoint", u.String()), zap.Error(err))
						time.Sleep(time.Second)
						continue
					}
					break
				}

				cli := &client{
					userId:        userId,
					token:         token,
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
			var msgTotal = config.SystemConfig.Mock.MessageCount
			ticker := time.NewTicker(time.Second)
			for {
				select {
				case <-ticker.C:
					{
						id := strconv.Itoa(rand.Intn(999) + 1)
						msg := &pb.Message{
							MsgId:      uuid.New().String(),
							MsgType:    types.TextMessage,
							Content:    time.Now().String(),
							FromType:   types.FromUser,
							FromId:     cli.userId,
							FromDevice: cli.deviceId,
							ToType:     types.ToUser,
							ToId:       "user-" + id,
							ToDevice:   "device-" + id,
						}
						body, _ := proto.Marshal(msg)
						err := cli.conn.WriteMessage(websocket.BinaryMessage, protocol.WebsocketCodec.Encode(protocol.Message, body))
						if err != nil {
							log.Warn("Error sending message", zap.Error(err))
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
