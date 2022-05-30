package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	c := context.Background()
	// go语言倾向于把变量名取短，因为在小代码块里面会很直观，变量名的长度会影响代码的可读性。
	mc, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017/coolcar?readPreference=primary&ssl=false"))
	if err != nil {
		panic(err)
	}
	col := mc.Database("coolcar").Collection("account")

	//insertRows(c, col)
	//findRow(c, col)
	//findManyRows(c, col)
	findRows(c, col)
}

// 查找多条数据，方法二
func findRows(c context.Context, col *mongo.Collection) {
	res, err := col.Find(c, bson.M{
		"open_id": "123",
	})
	if err != nil {
		panic(err)
	}
	// 遍历结果
	for res.Next(c) {
		var row struct {
			ID     primitive.ObjectID `bson:"_id"`
			OpenID string             `bson:"open_id"`
		}
		err = res.Decode(&row)
		if err != nil {
			panic(err)
		}
		fmt.Printf("%+v\n", row)
	}
}

// 查找多条数据，方法一
func findManyRows(c context.Context, col *mongo.Collection) {
	res, err := col.Find(c, bson.M{
		"open_id": "123",
	})
	if err != nil {
		panic(err)
	}
	var rows []struct {
		ID     primitive.ObjectID `bson:"_id"`
		OpenID string             `bson:"open_id"`
	}
	err = res.All(c, &rows)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", rows)
}

// 查找单条数据
func findRow(c context.Context, col *mongo.Collection) {
	// 查找第一条数据
	res := col.FindOne(c, bson.M{
		"open_id": "123",
	})
	fmt.Printf("%+v\n", res)
	var row struct {
		ID     primitive.ObjectID `bson:"_id"`
		OpenID string             `bson:"open_id"`
	}
	err := res.Decode(&row)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", row)
}

// 插入多条数据
func insertRows(c context.Context, col *mongo.Collection) {
	res, err := col.InsertMany(c, []interface{}{
		// bson是一个MongoDB内部的一种格式，二进制的数据格式，json很占空间，所以要把它序列化成二进制
		bson.M{
			"open_id": "123",
		},
		bson.M{
			"open_id": "456",
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", res)
}
