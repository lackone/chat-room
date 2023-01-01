package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Message struct {
	Id           string `bson:"_id"`
	UserIdentity string `bson:"user_identity"`
	RoomIdentity string `bson:"room_identity"`
	Data         string `bson:"data"`
	CreatedAt    int64  `bson:"created_at"`
	UpdatedAt    int64  `bson:"updated_at"`
}

// 获取集合名
func (this *Message) GetCollection() string {
	return "messages"
}

// 插入一条消息
func (this *Message) InsertOneMessage(m *Message) error {
	_, err := Mongodb.Collection(this.GetCollection()).InsertOne(context.Background(), m)
	return err
}

// 获取消息列表
func (this *Message) GetListByRoomId(roomIdentity string, skip, limit *int64) ([]*Message, error) {
	result := make([]*Message, 0)
	cur, err := Mongodb.Collection(this.GetCollection()).Find(context.Background(), bson.D{{"room_identity", roomIdentity}}, &options.FindOptions{
		Limit: limit,
		Skip:  skip,
		Sort:  bson.D{{"created_at", -1}},
	})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.Background()) {
		m := Message{}
		err := cur.Decode(&m)
		if err != nil {
			continue
		}
		result = append(result, &m)
	}
	return result, nil
}
