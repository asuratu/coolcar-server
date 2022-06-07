package main

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/server"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"log"
	"net/http"
)

func main() {
	c := context.Background()
	c, cancel := context.WithCancel(c)
	defer cancel()

	logger, err := server.NewZapLogger()
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

	servierConfig := []struct {
		name         string
		addr         string
		registerFunc func(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
	}{
		{
			name:         "auth",
			addr:         ":4001",
			registerFunc: authpb.RegisterAuthServiceHandlerFromEndpoint,
		},
		{
			name:         "rental",
			addr:         ":4002",
			registerFunc: rentalpb.RegisterTripServiceHandlerFromEndpoint,
		},
	}

	for _, config := range servierConfig {
		err := config.registerFunc(c, mux, config.addr, []grpc.DialOption{grpc.WithInsecure()})
		if err != nil {
			logger.Sugar().Fatalw("cannot start grpc gatway", "name", config.name, "addr", config.addr, "err", err)
		}
	}

	addr := ":6800"
	logger.Sugar().Infof("grpc gateway started at %s", addr)
	logger.Sugar().Fatal(http.ListenAndServe(addr, mux))
}
