package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRoom struct {
	Id           string `bson:"_id"`
	UserIdentity string `bson:"user_identity"`
	RoomIdentity string `bson:"room_identity"`
	RoomType     int8   `bson:"room_type"` //1群聊，2单聊
	CreatedAt    int64  `bson:"created_at"`
	UpdatedAt    int64  `bson:"updated_at"`
}

// 获取集合名
func (this *UserRoom) GetCollection() string {
	return "user_rooms"
}

// 获取用户房间
func (this *UserRoom) GetUserRoom(UserIdentity, RoomIdentity string) (*UserRoom, error) {
	ur := UserRoom{}
	err := Mongodb.Collection(this.GetCollection()).FindOne(context.Background(), bson.D{{"user_identity", UserIdentity}, {"room_identity", RoomIdentity}}).Decode(&ur)
	if err != nil {
		return nil, err
	}
	return &ur, nil
}

// 获取指定房间下所有用户
func (this *UserRoom) GetUsersByRooms(RoomIdentity string) ([]*UserRoom, error) {
	result := make([]*UserRoom, 0)

	cur, err := Mongodb.Collection(this.GetCollection()).Find(context.Background(), bson.D{{"room_identity", RoomIdentity}})
	if err != nil {
		return nil, err
	}

	for cur.Next(context.Background()) {
		ur := UserRoom{}
		err = cur.Decode(&ur)
		if err != nil {
			continue
		}
		result = append(result, &ur)
	}

	return result, nil
}

// 获取用户下所有单聊房间
func (this *UserRoom) GetSingleRoomsByUser(userIdentity string) ([]string, error) {
	result := make([]string, 0)

	cur, err := Mongodb.Collection(this.GetCollection()).Find(context.Background(), bson.D{{"user_identity", userIdentity}, {"room_type", 2}})
	if err != nil {
		return nil, err
	}
	for cur.Next(context.Background()) {
		ur := UserRoom{}
		err = cur.Decode(&ur)
		if err != nil {
			continue
		}
		result = append(result, ur.RoomIdentity)
	}
	return result, nil
}

// 获取是否为好友关系
func (this *UserRoom) IsFriend(userIdentity1, userIdentity2 string) (bool, error) {
	rooms, err := this.GetSingleRoomsByUser(userIdentity1)
	if err != nil {
		return false, err
	}
	cnt, err := Mongodb.Collection(this.GetCollection()).CountDocuments(context.Background(), bson.D{
		{"user_identity", userIdentity2},
		{"room_type", 2},
		{"room_identity", bson.D{{"$in", rooms}}},
	})
	if err != nil {
		return false, err
	}
	if cnt > 0 {
		return true, nil
	}
	return false, nil
}

// 添加记录
func (this *UserRoom) InsertOneUserRoom(ur *UserRoom) error {
	_, err := Mongodb.Collection(this.GetCollection()).InsertOne(context.Background(), ur)
	return err
}

// 获取单聊房间
func (this *UserRoom) GetSingleRoom(userIdentity1, userIdentity2 string) (string, error) {
	rooms, err := this.GetSingleRoomsByUser(userIdentity1)
	if err != nil {
		return "", err
	}
	ur := UserRoom{}
	err = Mongodb.Collection(this.GetCollection()).FindOne(context.Background(), bson.D{
		{"user_identity", userIdentity2},
		{"room_type", 2},
		{"room_identity", bson.D{{"$in", rooms}}},
	}).Decode(&ur)
	if err != nil {
		return "", err
	}
	return ur.RoomIdentity, nil
}

// 删除房间
func (this *UserRoom) DelByIdAndUser(roomIdentity string, userIdentity1, userIdentity2 string) error {
	_, err := Mongodb.Collection(this.GetCollection()).DeleteMany(context.Background(), bson.D{
		{"room_identity", roomIdentity},
		{"user_identity", bson.D{{"$in", []string{userIdentity1, userIdentity2}}}},
	})
	return err
}
