package auth

import (
	"context"
	"coolcar/shared/auth/token"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"io/ioutil"
	"os"
	"strings"
)

const (
	authorizationHeader = "authorization"
	bearerPrefix        = "Bearer "
)

// Interceptor creates a grpc auth interceptor.
func Interceptor(publicKeyFile string) (grpc.UnaryServerInterceptor, error) {
	f, err := os.Open(publicKeyFile)
	if err != nil {
		return nil, fmt.Errorf("cannnot open public key file: %v", err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("cannot read public key: %v", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(b)
	if err != nil {
		return nil, fmt.Errorf("canot parse public key: %v", err)
	}
	i := &interceptor{
		verifier: &token.JWTTokenVerifier{
			PublicKey: pubKey,
		},
	}
	return i.HandleRequest, nil
}

type tokenVerifier interface {
	Verify(token string) (string, error)
}

// share 文件夹里面文件不要直接暴露出去，所以这里用小写
type interceptor struct {
	// 结构前面加*，interface前面从不加*
	verifier tokenVerifier
}

func (i *interceptor) HandleRequest(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	tkn, err := tokenFromContext(ctx)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "no authorization token found in request")
	}
	aid, err := i.verifier.Verify(tkn)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization token: %v", err)
	}
	return handler(ContextWithAccountID(ctx, aid), req)
}

type accountIDKey struct {
	accountID string
}

// ContextWithAccountID 在ctx上面添加accountID
func ContextWithAccountID(ctx context.Context, aid string) context.Context {
	return context.WithValue(ctx, accountIDKey{}, aid)
}

// tokenFromContext 获取token
func tokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "no authorization token found in request")
	}
	tkn := ""
	for _, v := range md[authorizationHeader] {
		if strings.HasPrefix(v, bearerPrefix) {
			tkn = v[len(bearerPrefix):]
		}
	}
	if tkn == "" {
		return "", status.Errorf(codes.Unauthenticated, "no authorization token found in request")
	}
	return tkn, nil
}

// AccountIDFromContext get accountID from context
// Returns Unauthenticated error if no accountID found in context
func AccountIDFromContext(ctx context.Context) (string, error) {
	aid, ok := ctx.Value(accountIDKey{}).(string)
	if !ok {
		return "", status.Errorf(codes.Unauthenticated, "no accountID found in request")
	}
	return aid, nil
}
