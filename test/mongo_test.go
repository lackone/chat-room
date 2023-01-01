package test

import (
	"context"
	"github.com/lackone/chat-room/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"testing"
)

func TestFindOne(t *testing.T) {
	var u models.User
	models.Mongodb.Collection(u.GetCollection()).FindOne(context.Background(), bson.D{{}}).Decode(&u)
	t.Log(u)

	var m models.Message
	models.Mongodb.Collection(m.GetCollection()).FindOne(context.Background(), bson.D{{}}).Decode(&m)
	t.Log(m)

	var r models.Room
	models.Mongodb.Collection(r.GetCollection()).FindOne(context.Background(), bson.D{{}}).Decode(&r)
	t.Log(r)

	var urm models.UserRoom
	models.Mongodb.Collection(urm.GetCollection()).FindOne(context.Background(), bson.D{{}}).Decode(&urm)
	t.Log(urm)
}

func TestDelete(t *testing.T) {
	id := "63b19e2bf4a8ce396782f69f"
	objectId, err := primitive.ObjectIDFromHex(id)
	t.Log(objectId)
	t.Log(err)
	rr, err := models.Mongodb.Collection((&models.Room{}).GetCollection()).DeleteOne(context.Background(), bson.D{{"$or", bson.A{bson.D{{"_id", objectId}}, bson.D{{"_id", id}}}}})
	t.Log(rr.DeletedCount)
	t.Log(err)
}
