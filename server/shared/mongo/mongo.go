package mgutil

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

// Common field names
const (
	IDFieldName        = "_id"
	UpdatedAtFieldName = "updatedat"
)

type IdField struct {
	ID primitive.ObjectID `bson:"_id"`
}

type UpdatedAtField struct {
	UpdatedAt int64 `bson:"updatedat"`
}

// NewObjID returns a new ObjectID
var NewObjID = primitive.NewObjectID

// UpdatedAt returns a value suitable for UpdatedAt field
var UpdatedAt = func() int64 {
	// nano is 10^-9 second
	return time.Now().UnixNano()
}

// Set returns a $set update document
func Set(v interface{}) bson.M {
	return bson.M{"$set": v}
}

func SetOnInsert(v interface{}) bson.M {
	return bson.M{"$setOnInsert": v}
}
