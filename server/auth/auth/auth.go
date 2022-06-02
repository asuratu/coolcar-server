package auth

import (
	"context"
	authpb "coolcar/auth/api/gen/v1"
	"coolcar/auth/dao"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"time"
)

type Service struct {
	OpenIDResolver
	// Mongo 这里不建议设成接口，接口适合小的功能模块，比如这里的OpenIDResolver
	Mongo *dao.Mongo
	TokenGenerator
	TokenExpire time.Duration
	Logger      *zap.Logger
}

// OpenIDResolver resolvers an authorization code
// to an open id
type OpenIDResolver interface {
	Resolve(code string) (string, error)
}

// TokenGenerator generates a token for specified account
type TokenGenerator interface {
	GenerateToken(account string, expire time.Duration) (string, error)
}

func (s Service) Login(ctx context.Context, request *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	s.Logger.Info("received code", zap.String("code", request.Code))
	openID, err := s.OpenIDResolver.Resolve(request.Code)
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "cannot resolve openid:%v", err)
	}

	accountID, err := s.Mongo.ResolveAccountID(ctx, openID)
	if err != nil {
		s.Logger.Error("cannot resolve account id", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "cannot resolve account id:%v", err)
	}

	tkn, err := s.TokenGenerator.GenerateToken(accountID, s.TokenExpire)
	if err != nil {
		s.Logger.Error("cannot generate token", zap.Error(err))
		return nil, status.Errorf(codes.Internal, "cannot generate token:%v", err)
	}

	return &authpb.LoginResponse{
		//AccessToken: "Token for account id: " + accountID,
		AccessToken: tkn,
		ExpiresIn:   int32(s.TokenExpire.Seconds()),
	}, nil
}
