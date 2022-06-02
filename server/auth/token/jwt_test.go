package token

import (
	"github.com/dgrijalva/jwt-go"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		t.Fatalf("cannot parse private key: %v", err)
	}

	g := NewJWTTokenGen("coolcar/auth", key)

	// 固定每次的测试时间, 不用每次都计算时间，当前发布token的时间
	g.nowFuc = func() time.Time {
		return time.Unix(1654093996, 0)
	}

	tkn, err := g.GenerateToken("6294998159f23ab9ce8864a3", 2*time.Hour)
	if err != nil {
		t.Errorf("cannot generate token: %v", err)
	}

	want := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTQxMDExOTYsImlhdCI6MTY1NDA5Mzk5NiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjI5NDk5ODE1OWYyM2FiOWNlODg2NGEzIn0.S6k2hEHQwWIfqcygg3pyvrZ64kTjb46J-buJszM_cGmC6kdIY1vuqKziSFtzErZpz-Io8W5gpmfBYkCtGhCgP7IZpTqIltNd4kB6Pw7L-7m1QLqHHLX30LHiiPWmCzVOuk_L1nnYmMMIvEXAoSEp2ayACvSo7oLaQ_3Gc4bHgFx0cgLb1wPf0q1Axf2iROFt5CX4PtgXTGTgCK76ipZSc9RU6A8s1uRCUiYatqqK2QXn51Aq5WPdyGZWP871O8ATZmQuNJd346_PG3dAIRsluGfczwI8yL0UtNuMxeALed7gRJ-b9VAr677Nw929gIJkvsYJf5x4E7rgv3sfM5xY-w"

	if tkn != want {
		// %q is a quoted string
		t.Errorf("\ntoken = %q,\n want %q", tkn, want)
	}

}

const privateKey = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCbQB/AhqVNuxbT
12em8MrMclt+aUe7zYuqORBir5lS81SIwjjHR9RL9hqsLoO2x1i2ID9cf9nO8Lfw
FU7/Rvglmmti7Fhj2UyEekVQoFq2S9C1/gB6bS2z2y2wOeO6b1bYjwp4wLvqV4us
IWnMbZ2bKCgraa66WEjyUmUa394s6AmDOl8YVW03SG1I7P3dKo63mi/J8aNjDk/R
8IdQOMsuWzMGHmk8Hqg6ce+IdFRr6+BVD0mhcVQxqlN0jc3ssIjrLY4C3WoYjsp6
CwOrCIva5kKOsKxegWWnTxJ6kwSzt71iNp0tqryqnATZDvCGK26vmP0pDglqmkM+
GPMR4klrAgMBAAECggEAOARwMIik1qI9/1wG026ozhIpPzh/oJzu2xHR/rm7ifmw
s9PYptcdG/eF8kCqV+Yf9T83fYnILmofBGq74VJbMT5BpyT+U7DRci+oGQpzELnU
agZnZ8VDK1VXa/HHYLrRzDv4nE92vnyuMgKwaQnYR2a678cnO6elUoI2ZvcF9I+J
FbcUxlZuqzRYjhIPS5IAzR+BhTsxgMvm0qa50elAUA599rhJKNUvJoAqIpMMNkQK
8xn89QtnsPHvfqe1VrmkICNrFBnMs1EwclZXpG1Hy7QDhYnvJJAVHkinds+Hk4f1
HHUoFaeRUXj6dCxNJfW7at2JEi6Asvpz0OdAv5PHOQKBgQDKEv807ERmAiu8bbni
eMv/rO7hrvp5WPe9R3clgZ3WELdsScGd/yq8bZZDUOgSY8VoVfMKzQ2PsJfh6JGQ
FXPHZBdALR4P9QCSqD3ykEIGs2j9nCO+Qw3twVBAqbrAD/NLaekV7AfEI4GUIM+o
JeRuiTObm32NiBCMPZAPmane7QKBgQDErktPyZ/X5DL75Zi/VSffhSDeDJ/QG9ld
7pCARmosbsr0e9ZhgBHVTrNUG5qw8JRcLuhHzdCBei4TycQ6euSEBG8urCJQVFEL
H+E0BUasuRxMoLv42hp+ZpwJN7v0R0akLQsxNBCPyWigMIwWx6IsFflAsEmj+pSj
8dUYh3fmtwKBgQCqdeP6zOPV+TbTuOv5c1UC1OqeTnDUNIynisWjSffPQEK6gm1l
zn5KfVcoafOar7czEG35SoiKEbnNw9Ym6THFnVVPub+GTnKxRGMdXzuTU3zZkwFD
2mTBjzUXlxYGNm5Ry4HoEDds6VbBkfwaJ/zOkcaLGVuMLJ9o7fW8cy+s8QKBgB1z
1EdQdCVKQPBDw4nzYJMyRme6EDyDzxsn3G+5G7EnrjDKUqIrIYCF1ojj0Vhpzm23
gIUwJuccusWv0zjGqm2ylEuy7ziER9aYoYq+t4Sp+7jl4QA0+P6wsvEWbYj5G3T+
YcyudURy7r3+RxwqLPjZGYg8Xeq2XYfncPfJYYUlAoGBAJkDr8MpZQTUZChfTWjJ
fSLPiXb3iSKygUudn3c3OUhldURbUUUZmFnjWC9Q2n+oj/DBs+SmE0XQ0pZa03VM
T4d6oi0J+O3tpNOllKM9CX/38dNGKfYzmyFM+Wb9m03kixDGNnFD1TZ1Aiontrdq
sPdOkX5Rih331qjWERXdbSXt
-----END PRIVATE KEY-----
`
