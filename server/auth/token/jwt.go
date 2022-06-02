package token

import (
	"crypto/rsa"
	"github.com/dgrijalva/jwt-go"
	"time"
)

// JWTTokenGen 生成token的结构体
type JWTTokenGen struct {
	issuer     string
	nowFuc     func() time.Time
	privateKey *rsa.PrivateKey
}

// NewJWTTokenGen 相当于一个构造函数, 初始化一个JWTTokenGen
// issuer: 签发者
func NewJWTTokenGen(issuer string, privateKey *rsa.PrivateKey) *JWTTokenGen {
	return &JWTTokenGen{
		issuer:     issuer,
		nowFuc:     time.Now,
		privateKey: privateKey,
	}
}

func (t *JWTTokenGen) GenerateToken(accountID string, expire time.Duration) (string, error) {
	nowSec := t.nowFuc().Unix()
	// StandardClaims 不能增加自定义字段，MapClaims 可以
	tkn := jwt.NewWithClaims(jwt.SigningMethodRS512, jwt.StandardClaims{
		Issuer:    t.issuer,
		IssuedAt:  nowSec,
		ExpiresAt: nowSec + int64(expire.Seconds()),
		Subject:   accountID,
	})

	// 生成token
	return tkn.SignedString(t.privateKey)
}
