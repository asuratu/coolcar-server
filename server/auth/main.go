package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/auth"
	"coolcar/auth/dao"
	"coolcar/auth/token"
	"coolcar/auth/wechat"
	"coolcar/shared/server"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
)

type authConfig struct {
	mongoURI       string
	privateKeyFile string
	databaseName   string
	appID          string
	appSecret      string
}

func main() {
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf(" cannot create logger: %v\n", err)
	}

	authConfig := authConfig{
		mongoURI:       "mongodb://localhost:27017/coolcar?readPreference=primary&ssl=false",
		databaseName:   "coolcar",
		privateKeyFile: "/Users/asura/Code/go/ccmouse/coolcar/server/auth/private.key",
		appID:          "app_id",
		appSecret:      "app_secret",
	}

	c := context.Background()
	mongoClient, err := mongo.Connect(c, options.Client().ApplyURI(authConfig.mongoURI))
	if err != nil {
		logger.Fatal("cannot connect to mongo", zap.Error(err))
	}

	// 读取配置文件 private.key
	pkFile, err := os.Open(authConfig.privateKeyFile)
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

	config := &server.GRPCConfig{
		Name:              "auth",
		Addr:              ":4001",
		AuthPublicKeyFile: "",
		RegisterFunction: func(s *grpc.Server) {
			authpb.RegisterAuthServiceServer(s, &auth.Service{
				OpenIDResolver: &wechat.Service{
					AppID:     authConfig.appID,
					AppSecret: authConfig.appSecret,
				},
				Mongo:          dao.NewMongo(mongoClient.Database(authConfig.databaseName)),
				Logger:         logger,
				TokenExpire:    2 * time.Hour,
				TokenGenerator: token.NewJWTTokenGen("coolcar/auth", privateKey),
			})
		},
		Logger: logger,
	}

	logger.Sugar().Fatal(server.RunGRPCServer(config))
}
