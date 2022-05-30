package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net/http"
)

func main() {
	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer cancel()

	logger, err := newZapLogger()
	if err != nil {
		log.Fatalf(" cannot create logger: %v\n", err)
	}

	mux := runtime.NewServeMux(runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames:  true, // 使用proto字段名代替JSON中的骆驼式名称的字段名。
			UseEnumNumbers: true, // 使用protoc枚举定义的值作为数字发送。
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true, // 未知字段将被忽略
		},
	}))

	err = authpb.RegisterAuthServiceHandlerFromEndpoint(c, mux, ":4001", []grpc.DialOption{grpc.WithInsecure()})
	if err != nil {
		logger.Fatal("cannot start auth grpc gatway", zap.Error(err))
	}

	err = http.ListenAndServe(":6800", mux)
	if err != nil {
		logger.Fatal("cannot listen and server auth", zap.Error(err))
	}
}

// 自定义日志
func newZapLogger() (*zap.Logger, error) {
	cfg := zap.NewDevelopmentConfig()
	cfg.EncoderConfig.TimeKey = ""
	return cfg.Build()
}
