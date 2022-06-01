package dao

import (
	"context"
	mgo "coolcar/shared/mongo"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// openIDField 自然语言命名
const openIDField = "open_id"

// Mongo defines a mongo dao
type Mongo struct {
	// 不希望外面初始化 mongo.Collection
	col      *mongo.Collection
	newObjID func() primitive.ObjectID
}

// NewMongo returns a new mongo dao
func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col:      db.Collection("account"),
		newObjID: primitive.NewObjectID,
	}
}

// ResolveAccountID Resolve resolves an account id from open id
func (m *Mongo) ResolveAccountID(c context.Context, openID string) (string, error) {
	insertedID := m.newObjID()
	res := m.col.FindOneAndUpdate(c, bson.M{
		openIDField: openID,
	}, mgo.SetOnInsert(bson.M{
		mgo.IDField: insertedID,
		openIDField: openID,
	}), options.FindOneAndUpdate().
		SetUpsert(true).
		SetReturnDocument(options.After))

	if err := res.Err(); err != nil {
		return "", fmt.Errorf("cannot findOneAndUpdate: %v", err)
	}

	var row mgo.ObjID

	err := res.Decode(&row)

	if err != nil {
		return "", fmt.Errorf("cannot Decode reusult: %v", err)
	}

	return row.ID.Hex(), nil
}
