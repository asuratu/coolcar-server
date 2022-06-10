package dao

import (
	"context"
	mgutil "coolcar/shared/mongo"
	mongotesting "coolcar/shared/mongo/testing"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"testing"
)

var mongoURI string

func TestResolveAccountID(t *testing.T) {
	c := context.Background()
	mc, err := mongo.Connect(c, options.Client().ApplyURI(mongoURI))

	if err != nil {
		t.Fatalf("cannot connect to mongo: %v", err)
	}

	m := NewMongo(mc.Database("coolcar"))

	_, err = m.col.InsertMany(c, []interface{}{
		bson.M{
			mgutil.IDFieldName: mustObjectID("5f7c245ab0361e00ffb9fd6f"),
			openIDField:        "openid_1",
		},
		bson.M{
			mgutil.IDFieldName: mustObjectID("5f7c245ab0361e00ffb9fd70"),
			openIDField:        "openid_2",
		},
	})

	if err != nil {
		t.Fatalf("cannot insert many: %v", err)
	}

	mgutil.NewObjID = func() primitive.ObjectID {
		return mustObjectID("5f7c245ab0361e00ffb9fd71")
	}

	// 表格驱动测试
	cases := []struct {
		name   string
		openID string
		want   string
	}{
		{
			name:   "existing_user",
			openID: "openid_1",
			want:   "5f7c245ab0361e00ffb9fd6f",
		},
		{
			name:   "another_existing_user",
			openID: "openid_2",
			want:   "5f7c245ab0361e00ffb9fd70",
		},
		{
			name:   "new_user",
			openID: "openid_3",
			want:   "5f7c245ab0361e00ffb9fd71",
		},
	}

	for _, cc := range cases {
		t.Run(cc.name, func(t *testing.T) {
			got, err := m.ResolveAccountID(context.Background(), cc.openID)
			if err != nil {
				t.Errorf("cannot resolve account id: %v", err)
			}
			if got != cc.want {
				t.Fatalf("got %q, want %q", got, cc.want)
			}
		})
	}
}

func mustObjectID(hex string) primitive.ObjectID {
	objID, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		panic(err)
	}
	return objID
}

func TestMain(m *testing.M) {
	os.Exit(mongotesting.RunWithMongoInDocker(m, &mongoURI))
}
