package trip

import (
	"context"
	rentalpb "coolcar/rental/api/gen/v1"
	"coolcar/shared/auth"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	Logger *zap.Logger
}

func (s *Service) CreateTrip(c context.Context, req *rentalpb.CreateTripRequest) (*rentalpb.TripEntity, error) {
	// get account id from context
	aid, err := auth.AccountIDFromContext(c)
	if err != nil {
		return nil, err
	}
	s.Logger.Info("CreateTrip", zap.String("req", req.String()), zap.String("aid", aid.String()))
	return nil, status.Errorf(codes.Unimplemented, "not implemented")
}

func (s *Service) GetTrip(ctx context.Context, request *rentalpb.GetTripRequest) (*rentalpb.Trip, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) GetTrips(ctx context.Context, request *rentalpb.GetTripsRequest) (*rentalpb.GetTripsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) UpdateTrip(ctx context.Context, request *rentalpb.UpdateTripRequest) (*rentalpb.Trip, error) {
	//TODO implement me
	panic("implement me")
}
