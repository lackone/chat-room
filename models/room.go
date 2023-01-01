package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	Id           string `bson:"_id"`
	Number       string `bson:"number"`
	Name         string `bson:"name"`
	Info         string `bson:"info"`
	UserIdentity string `bson:"user_identity"`
	CreatedAt    int64  `bson:"created_at"`
	UpdatedAt    int64  `bson:"updated_at"`
}

// 获取集合名
func (this *Room) GetCollection() string {
	return "rooms"
}

// 添加记录
func (this *Room) InsertOneRoom(r *Room) error {
	_, err := Mongodb.Collection(this.GetCollection()).InsertOne(context.Background(), r)
	return err
}

// 删除记录
func (this *Room) DelOneRoom(id string) error {
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = Mongodb.Collection(this.GetCollection()).DeleteOne(context.Background(), bson.D{{"$or", bson.A{bson.D{{"_id", objectId}}, bson.D{{"_id", id}}}}})
	return err
}
