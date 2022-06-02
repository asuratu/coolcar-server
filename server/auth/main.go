package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/auth"
	"coolcar/auth/dao"
	"coolcar/auth/token"
	"coolcar/auth/wechat"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	logger, err := newZapLogger()
	if err != nil {
		log.Fatalf(" cannot create logger: %v\n", err)
	}

	listen, err := net.Listen("tcp", ":4001")
	if err != nil {
		logger.Fatal(" failed to listen", zap.Error(err))
	}

	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI("mongodb://localhost:27017/coolcar?readPreference=primary&ssl=false"))
	if err != nil {
		logger.Fatal("cannot connect to mongo", zap.Error(err))
	}

	// 读取配置文件 private.key
	pkFile, err := os.Open("/Users/asura/Code/go/ccmouse/coolcar/server/auth/private.key")
	if err != nil {
		logger.Fatal("cannot open private.key", zap.Error(err))
	}
	defer pkFile.Close()

	pkBytes, err := ioutil.ReadAll(pkFile)
	if err != nil {
		logger.Fatal("cannot read private.key", zap.Error(err))
	}

	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(pkBytes)
	if err != nil {
		logger.Fatal("cannot parse private.key", zap.Error(err))
	}

	s := grpc.NewServer()
	authpb.RegisterAuthServiceServer(s, &auth.Service{
		OpenIDResolver: &wechat.Service{
			AppID:     "app_id",
			AppSecret: "app_secret",
		},
		Mongo:          dao.NewMongo(mongoClient.Database("coolcar")),
		Logger:         logger,
		TokenExpire:    2 * time.Hour,
		TokenGenerator: token.NewJWTTokenGen("coolcar/auth", privateKey),
	})

	err = s.Serve(listen)
	logger.Fatal("cannot server", zap.Error(err))
}

// 自定义日志
func newZapLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.TimeKey = ""
	return cfg.Build()
}
