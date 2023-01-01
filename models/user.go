package models

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        string `bson:"_id"`
	Account   string `bson:"account"`
	Password  string `bson:"password"`
	Nickname  string `bson:"nickname"`
	Sex       int8   `bson:"sex"`
	Email     string `bson:"email"`
	Avatar    string `bson:"avatar"`
	CreatedAt int64  `bson:"created_at"`
	UpdatedAt int64  `bson:"updated_at"`
}

// 获取集合名
func (this *User) GetCollection() string {
	return "users"
}

// 获取用户
func (this *User) GetByAccountPassword(account, password string) (*User, error) {
	u := User{}
	err := Mongodb.Collection(this.GetCollection()).FindOne(context.Background(), bson.D{{"account", account}, {"password", password}}).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// 获取用户
func (this *User) GetByAccount(account string) (*User, error) {
	u := User{}
	err := Mongodb.Collection(this.GetCollection()).FindOne(context.Background(), bson.D{{"account", account}}).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// 获取用户
func (this *User) GetById(id string) (*User, error) {
	u := User{}
	//注意这里的id要转成ObjectID才能查询
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	err = Mongodb.Collection(this.GetCollection()).
		FindOne(context.Background(), bson.D{{"_id", objectId}}).
		Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// 获取用户
func (this *User) GetByEmail(email string) (*User, error) {
	u := User{}
	err := Mongodb.Collection(this.GetCollection()).FindOne(context.Background(), bson.D{{"email", email}}).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}

// 获取邮箱下用户数
func (this *User) GetCountByEmail(email string) (int64, error) {
	return Mongodb.Collection(this.GetCollection()).CountDocuments(context.Background(), bson.D{{"email", email}})
}

// 获取账号下用户数
func (this *User) GetCountByAccount(account string) (int64, error) {
	return Mongodb.Collection(this.GetCollection()).CountDocuments(context.Background(), bson.D{{"account", account}})
}

// 插入一个用户
func (this *User) InsertOneUser(u *User) error {
	_, err := Mongodb.Collection(this.GetCollection()).InsertOne(context.Background(), u)
	return err
}
