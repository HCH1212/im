package service

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"im/cache"
	"im/config"
	"im/pkg/e"
	"net/http"
	"strconv"
	"time"
)

const month = 60 * 60 * 24 * 30

type SendMsg struct { // 发送的消息
	Type    int    `json:"type"`
	Content string `json:"content"`
}

type ReplyMsg struct { // 回复消息
	From    string `json:"from"`
	Code    int    `json:"code"`
	Content string `json:"content"`
}

type Client struct { // 用户
	ID     string `json:"id"`
	SendID string `json:"send_id"`
	Socket *websocket.Conn
	Send   chan []byte
}

type Broadcast struct { // 广播i
	Client  *Client
	Message []byte
	Type    int
}

type ClientManager struct { // 用户管理
	Clients    map[string]*Client
	Broadcast  chan *Broadcast
	Reply      chan *ReplyMsg
	Register   chan *Client
	Unregister chan *Client
}

type Message struct { // 消息
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Content   string `json:"content,omitempty"`
}

var Manager = ClientManager{ // 用户管理实例
	Broadcast:  make(chan *Broadcast),
	Reply:      make(chan *ReplyMsg),
	Clients:    make(map[string]*Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func CreateID(uid, toUid string) string {
	return uid + "->" + toUid
}

func Handler(c *gin.Context) {
	uid := c.Query("uid")
	toUid := c.Query("toUid")
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(c.Writer, c.Request, nil) // 升级ws协议
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}

	// 创建一个用户实例
	client := &Client{
		ID:     CreateID(uid, toUid), // 1 -> 2
		SendID: CreateID(toUid, uid), // 2 -> 1
		Socket: conn,
		Send:   make(chan []byte),
	}

	// 将用户注册到用户管理上
	Manager.Register <- client
	go client.Read()
	go client.Write()
}

func (manager *Client) Read() {
	defer func() {
		Manager.Unregister <- manager
		_ = manager.Socket.Close()
	}()
	for {
		// 处理 WebSocket 心跳
		manager.Socket.PongHandler()

		sendMsg := &SendMsg{}
		err := manager.Socket.ReadJSON(sendMsg)
		if err != nil {
			logrus.Println("数据格式不正确：", err.Error())
			break
		}

		if sendMsg.Type == 1 { // 1 -> 2 发送消息
			r1, _ := cache.RedisClient.Get(context.Background(), manager.ID).Result()
			r2, _ := cache.RedisClient.Get(context.Background(), manager.SendID).Result()
			if r1 > "3" && r2 == "" {
				replyMsg := &ReplyMsg{ // 1给2发了三条消息，2没有回，就停止1发送
					Code:    e.WebsocketLimit,
					Content: e.GetMsg(e.WebsocketLimit),
				}
				msg, _ := json.Marshal(replyMsg)
				_ = manager.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			}

			cache.RedisClient.Incr(context.Background(), manager.ID)
			_, _ = cache.RedisClient.Expire(context.Background(), manager.ID, time.Hour*24*30*3).Result() // 防止过快分手，建立连接三个月过期

			Manager.Broadcast <- &Broadcast{
				Client:  manager,
				Message: []byte(sendMsg.Content), // 发送过来的消息
			}
		} else if sendMsg.Type == 2 {
			// 获取历史消息
			timeT, err := strconv.Atoi(sendMsg.Content)
			if err != nil {
				timeT = 999999
			}
			results, _ := FindMany(config.MongoDBName, manager.SendID, manager.ID, int64(timeT), 10) // 获取十条历史消息
			if len(results) > 10 {
				results = results[:10]
			} else if len(results) == 0 {
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "到底了",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = manager.Socket.WriteMessage(websocket.TextMessage, msg)
				continue
			}
			for _, result := range results {
				replyMsg := &ReplyMsg{
					From:    result.From,
					Content: result.Msg,
				}
				msg, _ := json.Marshal(replyMsg)
				_ = manager.Socket.WriteMessage(websocket.TextMessage, msg)
			}
		}
	}
}

func (manager *Client) Write() {
	defer func() {
		_ = manager.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-manager.Send:
			if !ok {
				_ = manager.Socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			replyMsg := &ReplyMsg{
				Code:    e.WebsocketSuccessMessage,
				Content: e.GetMsg(e.WebsocketSuccessMessage) + string(message),
			}
			msg, _ := json.Marshal(replyMsg)
			_ = manager.Socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
