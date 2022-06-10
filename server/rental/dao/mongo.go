package dao

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	mgutil "coolcar/shared/mongo"
	"go.mongodb.org/mongo-driver/mongo"
)

// Mongo defines a mongo dao
type Mongo struct {
	// 不希望外面初始化 mongo.Collection
	col *mongo.Collection
}

// NewMongo returns a new mongo dao
func NewMongo(db *mongo.Database) *Mongo {
	return &Mongo{
		col: db.Collection("trip"),
	}
}

// TripRecord 定义Trip的表数据结构
type TripRecord struct {
	mgutil.IdField
	mgutil.UpdatedAtField
	Trip *rentalpb.Trip `bson:"trip"`
}

func (m *Mongo) CreateTrip(c context.Context, trip *rentalpb.Trip) (*TripRecord, error) {
	r := &TripRecord{
		Trip: trip,
	}
	r.ID = mgutil.NewObjID()
	r.UpdatedAt = mgutil.UpdatedAt()
	_, err := m.col.InsertOne(c, r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
