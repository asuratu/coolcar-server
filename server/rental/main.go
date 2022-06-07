package main

import (
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/rental/trip"
	"coolcar/shared/server"
	"google.golang.org/grpc"
	"log"
)

func main() {
	logger, err := server.NewZapLogger()
	if err != nil {
		log.Fatalf(" cannot create logger: %v\n", err)
	}

	config := &server.GRPCConfig{
		Name:              "trip",
		Addr:              ":4002",
		AuthPublicKeyFile: "shared/auth/public.key",
		RegisterFunction: func(s *grpc.Server) {
			rentalpb.RegisterTripServiceServer(s, &trip.Service{
				Logger: logger,
			})
		},
		Logger: logger,
	}

	logger.Sugar().Fatal(server.RunGRPCServer(config))
}
