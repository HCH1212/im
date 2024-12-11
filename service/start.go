package service

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"im/config"
	"im/pkg/e"
	"log"
)

func (manager *ClientManager) Start() {
	for {
		fmt.Println("------监听管道通信------")
		select {
		case conn := <-Manager.Register:
			fmt.Printf("有新连接: %v", conn.ID)
			Manager.Clients[conn.ID] = conn // 把连接放到用户管理上
			replyMsg := ReplyMsg{
				Code:    e.WebsocketSuccess,
				Content: "已经连接服务器",
			}
			msg, _ := json.Marshal(replyMsg)
			_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
		case conn := <-Manager.Unregister:
			fmt.Println("连接失败" + conn.ID)
			if _, ok := Manager.Clients[conn.ID]; ok {
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketEnd,
					Content: "连接中断",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = conn.Socket.WriteMessage(websocket.TextMessage, msg)
				close(conn.Send)
				delete(Manager.Clients, conn.ID)
			}
		case broadcast := <-Manager.Broadcast: // if 1->2
			message := broadcast.Message
			sendID := broadcast.Client.SendID // 2->1
			flag := false                     // 默认对方不在线
			for id, conn := range Manager.Clients {
				if id != sendID {
					continue
				}
				select {
				case conn.Send <- message:
					flag = true
				default:
					close(conn.Send)
					delete(Manager.Clients, conn.ID)
				}
			}
			id := broadcast.Client.ID // 1->2
			if flag {
				replyMsg := &ReplyMsg{
					Code:    e.WebsocketOnlineReply,
					Content: "对方在线应答",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				err := InsertMsg(config.MongoDBName, id, string(message), 1, int64(3*month)) // 1 表示已读(在线)
				if err != nil {
					fmt.Println("insertOne err:", err)
				}
			} else {
				log.Println("对方不在线")
				replyMsg := ReplyMsg{
					Code:    e.WebsocketOfflineReply,
					Content: "对方不在线应答",
				}
				msg, err := json.Marshal(replyMsg)
				_ = broadcast.Client.Socket.WriteMessage(websocket.TextMessage, msg)
				err = InsertMsg(config.MongoDBName, id, string(message), 0, int64(3*month))
				if err != nil {
					fmt.Println("InsertOneMsg Err", err)
				}
			}
		}
	}
}
