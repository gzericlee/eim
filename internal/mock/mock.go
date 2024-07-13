package mock

import (
	"context"
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
	"go.uber.org/atomic"
	"go.uber.org/zap"

	"eim/internal/config"
	"eim/internal/gateway/protocol"
	"eim/internal/model"
	"eim/internal/model/consts"
	seqrpc "eim/internal/seq/rpc"
	"eim/pkg/httputil"
	"eim/pkg/log"
)

type Server struct {
	mockUserMsgCount  int
	mockGroupMsgCount int
	mockClientCount   int
	mockStartUserId   int
	mockStartGroupId  int
	mockSendCount     int
	connectedCount    *atomic.Int64
	sendCount         *atomic.Int64
	ackCount          *atomic.Int64
	msgCount          *atomic.Int64
	invalidCount      *atomic.Int64
	hbCount           *atomic.Int64
	clients           sync.Map
	seqRpc            *seqrpc.Client
	connectedConsume  time.Duration
	connectionStart   time.Time
	sendStart         time.Time
}

func NewMockServer(seqRpc *seqrpc.Client, userMsgCount, groupMsgCount, clientCount, startUserId, startGroupId, sendCount int) *Server {
	return &Server{
		mockUserMsgCount:  userMsgCount,
		mockGroupMsgCount: groupMsgCount,
		mockClientCount:   clientCount,
		mockStartGroupId:  startGroupId,
		mockStartUserId:   startUserId,
		mockSendCount:     sendCount,
		connectedCount:    atomic.NewInt64(0),
		sendCount:         atomic.NewInt64(0),
		ackCount:          atomic.NewInt64(0),
		msgCount:          atomic.NewInt64(0),
		invalidCount:      atomic.NewInt64(0),
		hbCount:           atomic.NewInt64(0),
		seqRpc:            seqRpc,
	}
}

func (its *Server) Start() {
	logging.SetLevel(logging.LevelNone)

	engine := nbhttp.NewEngine(nbhttp.Config{})
	err := engine.Start()
	if err != nil {
		log.Error("nbio.Start failed: %v\n", zap.Error(err))
		return
	}

	wg := sync.WaitGroup{}

	concurrency := make(chan struct{}, 500)
	its.connectionStart = time.Now()

	gateways, err := httputil.DoRequest[[]*model.Gateway](context.Background(), "http://127.0.0.1:10060/gateways", http.MethodGet, http.Header{
		"Authorization": []string{"Basic " + base64.StdEncoding.EncodeToString([]byte("user-1@bingo:"+"pass@word1"))},
	}, nil, true)
	if err != nil {
		panic(err)
	}

	go its.printStats()

	for i := its.mockStartUserId; i < its.mockStartUserId+its.mockClientCount; i++ {
		wg.Add(1)
		concurrency <- struct{}{}
		go func(i int) {
			defer func() {
				wg.Done()
				<-concurrency
			}()

			gateway := gateways[i%len(gateways)]

			id := strconv.Itoa(i)
			u := url.URL{Scheme: "ws", Host: fmt.Sprintf("%s:%d", gateway.Ip, gateway.Port), Path: "/"}

			userId := "user-" + id
			deviceId := "device-" + id
			deviceName := "linux-" + id
			deviceVersion := "1.0.0"
			deviceType := consts.LinuxDevice

			auth := base64.StdEncoding.EncodeToString([]byte(userId + "@bingo:" + "pass@word1"))
			token := "Basic " + auth

			client := &Client{
				userId:        userId,
				token:         token,
				deviceId:      deviceId,
				deviceName:    deviceName,
				deviceType:    deviceType,
				deviceVersion: deviceVersion,
				connected:     make(chan bool),
			}

			dialer := &websocket.Dialer{
				Engine:   engine,
				Upgrader: its.upgrade(client),
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			}

			for {
				client.conn, _, err = dialer.Dial(u.String(), http.Header{
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

			client.conn.SetSession(client)

			its.clients.Store(client.userId, client)

			<-client.connected

			its.connectedCount.Add(1)
		}(i)
		if err != nil {
			log.Error("Error submitting task", zap.Error(err))
			return
		}
	}

	wg.Wait()

	its.connectedConsume = time.Since(its.connectionStart)
	its.sendStart = time.Now()

	go its.mockHeartbeat()

	go its.mockUserMessage()

	go its.mockGroupMessage()
}

func (its *Server) mockHeartbeat() {
	go func() {
		ticker := time.NewTicker(time.Second * 30)
		for {
			select {
			case <-ticker.C:
				{
					its.clients.Range(func(key, value interface{}) bool {
						client := value.(*Client)
						_ = client.conn.WriteMessage(websocket.PingMessage, []byte(time.Now().String()))
						return true
					})
				}
			}
		}
	}()
}

func (its *Server) mockUserMessage() {
	if its.mockUserMsgCount <= 0 {
		return
	}
	for i := its.mockSendCount; i > 0; i-- {
		its.clients.Range(func(key, value interface{}) bool {
			go func(client *Client) {
				var msgTotal = its.mockUserMsgCount
				for {
					id := rand.Int31n(int32(config.SystemConfig.Mock.ClientCount)) + 1
					toId := fmt.Sprintf("user-%d", id)
					if client.userId == toId {
						continue
					}
					msgId, err := its.seqRpc.SnowflakeId()
					if err != nil {
						log.Error("Error getting snowflake id: %v，%v", zap.String("id", client.userId), zap.Error(err))
						continue
					}
					msg := &model.Message{
						MsgId:      msgId,
						MsgType:    consts.TextMessage,
						Content:    time.Now().String(),
						FromType:   consts.FromUser,
						FromId:     client.userId,
						FromDevice: client.deviceId,
						FromTenant: "bingo",
						ToType:     consts.ToUser,
						ToId:       toId,
						ToTenant:   "bingo",
					}
					its.sendMessage(client, msg)
					msgTotal--
					if msgTotal == 0 {
						return
					}
				}
			}(value.(*Client))
			return true
		})
		time.Sleep(time.Second * 30)
	}
}

func (its *Server) mockGroupMessage() {
	if its.mockGroupMsgCount <= 0 {
		return
	}
	for i := its.mockSendCount; i > 0; i-- {
		its.clients.Range(func(key, value interface{}) bool {
			go func(client *Client) {
				var msgTotal = its.mockGroupMsgCount
				for {
					id := rand.Int31n(int32(its.mockClientCount)) + 1
					toId := fmt.Sprintf("group-%d", id)
					if client.userId == toId {
						continue
					}
					msgId, err := its.seqRpc.SnowflakeId()
					if err != nil {
						log.Error("Error getting snowflake id: %v，%v", zap.String("id", client.userId), zap.Error(err))
						continue
					}
					msg := &model.Message{
						MsgId:      msgId,
						MsgType:    consts.TextMessage,
						Content:    time.Now().String(),
						FromType:   consts.FromUser,
						FromId:     client.userId,
						FromDevice: client.deviceId,
						FromTenant: "bingo",
						ToType:     consts.ToGroup,
						ToId:       toId,
						ToTenant:   "bingo",
					}
					its.sendMessage(client, msg)
					msgTotal--
					if msgTotal == 0 {
						return
					}
				}
			}(value.(*Client))
			return true
		})
		time.Sleep(time.Second * 30)
	}
}

func (its *Server) printStats() {
	var ticker = time.NewTicker(time.Second)
	var so, ro, sl, rl int64
	lastSend := &atomic.Int64{}
	lastReceive := &atomic.Int64{}

	for {
		select {
		case <-ticker.C:
			{
				sc := its.sendCount.Load()
				if sc > 0 && lastSend.Load() < sc {
					sl++
					so = sc / sl
				}
				mc := its.msgCount.Load()
				if mc > 0 && lastReceive.Load() < mc {
					rl++
					ro = mc / rl
				}

				lastSend.Store(sc)
				lastReceive.Store(mc)

				t := table.NewWriter()
				t.SetOutputMirror(os.Stdout)
				t.AppendHeader(table.Row{"devices", "Send", "Send Tps", "Received", "Invalid", "Received Tps", "Ack", "Heartbeat", "Goroutines"})
				t.AppendRows([]table.Row{{
					its.connectedCount.Load(), its.sendCount.Load(), so, its.msgCount.Load(), its.invalidCount.Load(), ro, its.ackCount.Load(), its.hbCount.Load(), runtime.NumGoroutine(),
				}})
				t.AppendSeparator()
				t.Render()

				if its.mockUserMsgCount > 0 {
					if its.sendCount.Load() == its.msgCount.Load() && its.sendCount.Load() == int64(its.mockClientCount*its.mockUserMsgCount) && its.msgCount.Load() != 0 {
						log.Info(fmt.Sprintf("Mock completed，Connection %v : %v，Send %v : %v",
							its.mockClientCount,
							its.connectedConsume,
							its.mockClientCount*its.mockUserMsgCount,
							time.Since(its.sendStart)))
						return
					}
				}
			}
		}
	}
}

func (its *Server) sendMessage(client *Client, msg *model.Message) {
	body, err := proto.Marshal(msg)
	if err != nil {
		log.Warn("Error marshalling message", zap.Error(err))
		return
	}
	err = client.conn.WriteMessage(websocket.BinaryMessage, protocol.WebsocketCodec.Encode(protocol.Message, body))
	if err != nil {
		log.Warn("Error sending message", zap.Error(err))
	} else {
		its.sendCount.Add(1)
	}
}
