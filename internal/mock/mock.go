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
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/lesismal/nbio/extension/tls"
	"github.com/lesismal/nbio/logging"
	"github.com/lesismal/nbio/nbhttp"
	"github.com/lesismal/nbio/nbhttp/websocket"
	"github.com/panjf2000/ants"
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"eim/internal/config"
	"eim/internal/gateway/protocol"
	"eim/internal/model"
	seqrpc "eim/internal/seq/rpc"
	"eim/util/log"
)

var connectedCount = &atomic.Int64{}
var sendCount = &atomic.Int64{}
var ackCount = &atomic.Int64{}
var msgCount = &atomic.Int64{}
var invalidCount = &atomic.Int64{}
var hbCount = &atomic.Int64{}
var conns sync.Map
var connectedConsume time.Duration

func Do() {
	var connectionStart, sendStart time.Time
	logging.SetLevel(logging.LevelNone)

	seqRpc, err := seqrpc.NewClient(config.SystemConfig.Etcd.Endpoints.Value())
	if err != nil {
		log.Error("Error creating seq rpc client", zap.Error(err))
		return
	}

	go func() {
		var ticker = time.NewTicker(time.Second)
		var so, ro, sl, rl int64
		lastSent := &atomic.Int64{}
		lastRecv := &atomic.Int64{}

		for {
			select {
			case <-ticker.C:
				{
					sc := sendCount.Load()
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
					t.AppendHeader(table.Row{"devices", "Send", "Send Tps", "Received", "Invalid", "Received Tps", "Ack", "Heartbeat", "Goroutines"})
					t.AppendRows([]table.Row{{
						connectedCount.Load(), sendCount.Load(), so, msgCount.Load(), invalidCount.Load(), ro, ackCount.Load(), hbCount.Load(), runtime.NumGoroutine(),
					}})
					t.AppendSeparator()
					t.Render()

					if config.SystemConfig.Mock.MessageCount > 0 {
						if sendCount.Load() == msgCount.Load() && sendCount.Load() == int64(config.SystemConfig.Mock.ClientCount*config.SystemConfig.Mock.MessageCount) && msgCount.Load() != 0 {
							log.Info(fmt.Sprintf("Mock completed，Connection %v : %v，Send %v : %v",
								config.SystemConfig.Mock.ClientCount,
								connectedConsume,
								config.SystemConfig.Mock.ClientCount*config.SystemConfig.Mock.MessageCount,
								time.Since(sendStart)))
							return
						}
					} else {
						//if connectedCount.Load() == int64(config.SystemConfig.Mock.ClientCount) {
						//	log.Info(fmt.Sprintf("Mock completed，Connection %v : %v",
						//		config.SystemConfig.Mock.ClientCount,
						//		connectedConsume))
						//	return
						//}
					}
				}
			}
		}
	}()

	engine := nbhttp.NewEngine(nbhttp.Config{})
	err = engine.Start()
	if err != nil {
		log.Error("nbio.Start failed: %v\n", zap.Error(err))
		return
	}

	wg := sync.WaitGroup{}
	pool, err := ants.NewPoolPreMalloc(1024)
	if err != nil {
		log.Error("ants.NewPoolPreMalloc failed: %v\n", zap.Error(err))
		return
	}

	connectionStart = time.Now()

	for i := 1; i <= config.SystemConfig.Mock.ClientCount; i++ {
		wg.Add(1)
		err := pool.Submit(func(i int) func() {
			return func() {
				defer wg.Done()
				id := strconv.Itoa(i)
				u := url.URL{Scheme: "ws", Host: config.SystemConfig.Mock.EimEndpoints.Value()[i%len(config.SystemConfig.Mock.EimEndpoints.Value())], Path: "/"}

				userId := "user-" + id
				deviceId := "device-" + id
				deviceName := "linux-" + id
				deviceVersion := "1.0.0"
				deviceType := model.LinuxDevice

				dialer := &websocket.Dialer{
					Engine:   engine,
					Upgrader: newUpgrader(),
					TLSClientConfig: &tls.Config{
						InsecureSkipVerify: true,
					},
				}

				auth := base64.StdEncoding.EncodeToString([]byte(userId + "@bingo:" + "pass@word1"))
				token := "Basic " + auth

				var conn *websocket.Conn
				for {
					conn, _, err = dialer.Dial(u.String(), http.Header{
						"Authorization": []string{token},
						"deviceId":      []string{deviceId},
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
		if err != nil {
			log.Error("Error submitting task", zap.Error(err))
			return
		}
	}

	wg.Wait()

	connectedConsume = time.Since(connectionStart)

	sendStart = time.Now()

	conns.Range(func(key, value interface{}) bool {
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

	if config.SystemConfig.Mock.MessageCount > 0 {
		conns.Range(func(key, value interface{}) bool {
			go func(cli *client) {
				var msgTotal = config.SystemConfig.Mock.MessageCount
				for {
					id := rand.Int31n(int32(config.SystemConfig.Mock.ClientCount)) + 1
					toId := fmt.Sprintf("user-%d", id)
					toDevice := fmt.Sprintf("device-%d", id)
					if cli.userId == toId {
						continue
					}
					msgId, err := seqRpc.SnowflakeId()
					if err != nil {
						log.Error("Error getting snowflake id: %v，%v", zap.String("id", cli.userId), zap.Error(err))
						continue
					}
					msg := &model.Message{
						MsgId:      msgId,
						MsgType:    model.TextMessage,
						Content:    time.Now().String(),
						FromType:   model.FromUser,
						FromId:     cli.userId,
						FromDevice: cli.deviceId,
						ToType:     model.ToUser,
						ToId:       toId,
						ToDevice:   toDevice,
					}
					body, _ := proto.Marshal(msg)
					err = cli.conn.WriteMessage(websocket.BinaryMessage, protocol.WebsocketCodec.Encode(protocol.Message, body))
					if err != nil {
						log.Warn("Error sending message", zap.Error(err))
						return
					}
					sendCount.Add(1)
					msgTotal--
					if msgTotal == 0 {
						return
					}
				}
			}(value.(*client))
			return true
		})
	}
}
