package dao

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestResolveAccountID(t *testing.T) {
	c := context.Background()
	mc, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017/coolcar?readPreference=primary&ssl=false"))

	if err != nil {
		t.Fatalf("cannot connect to mongo: %v", err)
	}

	m := NewMongo(mc.Database("coolcar"))

	id, err := m.ResolveAccountID(c, "123")

	if err != nil {
		t.Errorf("cannot resolve account id for 123: %v", err)
	} else {
		want := "629482ec61fc41116c3b89b1"
		if id != want {
			// 如果id不等于want，则打印错误信息
			// %q 单引号，%v 值，%T 类型 %s 字符串
			t.Errorf("got %q, want %q", id, want)
		}
	}
}
