package service

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"im/config"
	"im/model/ws"
	"sort"
	"time"
)

type SendSortMsg struct {
	Content  string `json:"content"`
	Read     int64  `json:"read"`
	CreateAt int64  `json:"create_at"`
}

func InsertMsg(database, id, content string, read int64, expire int64) error {
	// 插入到MongoDB
	collection := config.MongoDBClient.Database(database).Collection(id) // 没有id集合的话，会自动创建
	comment := ws.Trainer{
		Content:   content,
		Read:      read,
		StartTime: time.Now().Unix(),
		EndTime:   time.Now().Unix() + expire,
	}
	_, err := collection.InsertOne(context.Background(), comment)
	return err
}

func FindMany(database, sendID, id string, time int64, pageSize int) (results []ws.Result, err error) {
	var resultMe []ws.Trainer  // id
	var resultYou []ws.Trainer // sendID

	sendIDCollection := config.MongoDBClient.Database(database).Collection(sendID)
	idCollection := config.MongoDBClient.Database(database).Collection(id)

	// 必须加bson.D{}，不然查询参数少一个，查询为空
	sendIdTimeCursor, err := sendIDCollection.Find(context.Background(), bson.D{},
		options.Find().SetSort(bson.D{{"startTime", -1}}), options.Find().SetLimit(int64(pageSize)))
	idTimeCursor, err := idCollection.Find(context.Background(), bson.D{},
		options.Find().SetSort(bson.D{{"startTime", -1}}), options.Find().SetLimit(int64(pageSize)))

	err = sendIdTimeCursor.All(context.Background(), &resultYou) // sendId 对面发过来的
	err = idTimeCursor.All(context.Background(), &resultMe)      // Id 发给对面的

	results, _ = AppendAndSort(resultMe, resultYou)

	return
}

func AppendAndSort(resultsMe, resultsYou []ws.Trainer) (results []ws.Result, err error) {
	for _, r := range resultsMe {
		sendSort := SendSortMsg{ // 构造返回msg
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := ws.Result{ // 构造返回所有的内容，包括传送者
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "me",
		}
		results = append(results, result)
	}
	for _, r := range resultsYou {
		sendSort := SendSortMsg{
			Content:  r.Content,
			Read:     r.Read,
			CreateAt: r.StartTime,
		}
		result := ws.Result{
			StartTime: r.StartTime,
			Msg:       fmt.Sprintf("%v", sendSort),
			From:      "you",
		}
		results = append(results, result)
	}
	// 最后进行排序
	sort.Slice(results, func(i, j int) bool { return results[i].StartTime < results[j].StartTime })
	return results, nil
}
