package server

import (
	"coolcar/shared/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"net"
)

// GRPCConfig is gRPC server configuration.
type GRPCConfig struct {
	Name              string
	Addr              string
	AuthPublicKeyFile string
	RegisterFunction  func(s *grpc.Server)
	Logger            *zap.Logger
}

// RunGRPCServer runs gRPC server and returns error if any.
func RunGRPCServer(c *GRPCConfig) error {
	// 日志命名
	nameField := zap.String("name", c.Name)
	listen, err := net.Listen("tcp", c.Addr)
	if err != nil {
		c.Logger.Fatal(" failed to listen", nameField, zap.Error(err))
	}

	var opts []grpc.ServerOption

	if c.AuthPublicKeyFile != "" {
		in, err := auth.Interceptor(c.AuthPublicKeyFile)
		if err != nil {
			c.Logger.Fatal(" failed to create interceptor", nameField, zap.Error(err))
		}
		opts = append(opts, grpc.UnaryInterceptor(in))
	}

	s := grpc.NewServer(opts...)

	c.RegisterFunction(s)

	return s.Serve(listen)
}
