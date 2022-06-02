package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAm0AfwIalTbsW09dnpvDK
zHJbfmlHu82LqjkQYq+ZUvNUiMI4x0fUS/YarC6DtsdYtiA/XH/ZzvC38BVO/0b4
JZprYuxYY9lMhHpFUKBatkvQtf4Aem0ts9stsDnjum9W2I8KeMC76leLrCFpzG2d
mygoK2muulhI8lJlGt/eLOgJgzpfGFVtN0htSOz93SqOt5ovyfGjYw5P0fCHUDjL
LlszBh5pPB6oOnHviHRUa+vgVQ9JoXFUMapTdI3N7LCI6y2OAt1qGI7KegsDqwiL
2uZCjrCsXoFlp08SepMEs7e9YjadLaq8qpwE2Q7whitur5j9KQ4JappDPhjzEeJJ
awIDAQAB
-----END PUBLIC KEY-----`

func TestVerify(t *testing.T) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		t.Fatalf("cannot parse public key: %v", err)
	}

	v := &JWTTokenVerifier{
		PublicKey: pubKey,
	}

	cases := []struct {
		name    string
		tkn     string
		now     time.Time
		want    string
		wantErr bool
	}{
		{
			name: "valid_token",
			tkn:  "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTQxMDExOTYsImlhdCI6MTY1NDA5Mzk5NiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjI5NDk5ODE1OWYyM2FiOWNlODg2NGEzIn0.S6k2hEHQwWIfqcygg3pyvrZ64kTjb46J-buJszM_cGmC6kdIY1vuqKziSFtzErZpz-Io8W5gpmfBYkCtGhCgP7IZpTqIltNd4kB6Pw7L-7m1QLqHHLX30LHiiPWmCzVOuk_L1nnYmMMIvEXAoSEp2ayACvSo7oLaQ_3Gc4bHgFx0cgLb1wPf0q1Axf2iROFt5CX4PtgXTGTgCK76ipZSc9RU6A8s1uRCUiYatqqK2QXn51Aq5WPdyGZWP871O8ATZmQuNJd346_PG3dAIRsluGfczwI8yL0UtNuMxeALed7gRJ-b9VAr677Nw929gIJkvsYJf5x4E7rgv3sfM5xY-w",
			now:  time.Unix(1654094000, 0),
			want: "6294998159f23ab9ce8864a3",
		},
		{
			name:    "token_expired",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTQxMDExOTYsImlhdCI6MTY1NDA5Mzk5NiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjI5NDk5ODE1OWYyM2FiOWNlODg2NGE0In0.S6k2hEHQwWIfqcygg3pyvrZ64kTjb46J-buJszM_cGmC6kdIY1vuqKziSFtzErZpz-Io8W5gpmfBYkCtGhCgP7IZpTqIltNd4kB6Pw7L-7m1QLqHHLX30LHiiPWmCzVOuk_L1nnYmMMIvEXAoSEp2ayACvSo7oLaQ_3Gc4bHgFx0cgLb1wPf0q1Axf2iROFt5CX4PtgXTGTgCK76ipZSc9RU6A8s1uRCUiYatqqK2QXn51Aq5WPdyGZWP871O8ATZmQuNJd346_PG3dAIRsluGfczwI8yL0UtNuMxeALed7gRJ-b9VAr677Nw929gIJkvsYJf5x4E7rgv3sfM5xY-w",
			now:     time.Unix(1654194000, 0),
			wantErr: true,
		},
		{
			name:    "bad_token",
			tkn:     "bad_token",
			now:     time.Unix(1517239122, 0),
			wantErr: true,
		},
		{
			name:    "wrong_signature",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTQxMDExOTYsImlhdCI6MTY1NDA5Mzk5NiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjI5NDk5ODE1OWYyM2FiOWNlODg2NGEzIn0.S6k2hEHQwWIfqcygg3pyvrZ64kTjb46J-buJszM_cGmC6kdIY1vuqKziSFtzErZpz-Io8W5gpmfBYkCtGhCgP7IZpTqIltNd4kB6Pw7L-7m1QLqHHLX30LHiiPWmCzVOuk_L1nnYmMMIvEXAoSEp2ayACvSo7oLaQ_3Gc4bHgFx0cgLb1wPf0q1Axf2iROFt5CX4PtgXTGTgCK76ipZSc9RU6A8s1uRCUiYatqqK2QXn51Aq5WPdyGZWP871O8ATZmQuNJd346_PG3dAIRsluGfczwI8yL0UtNuMxeALed7gRJ-b9VAr677Nw929gIJkvsYJf5x4E7rgv3sfM5xY-w",
			now:     time.Unix(1516239122, 0),
			wantErr: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			jwt.TimeFunc = func() time.Time {
				return c.now
			}
			accountID, err := v.Verify(c.tkn)

			if !c.wantErr && err != nil {
				t.Errorf("verification failed: %v", err)
			}

			if c.wantErr && err == nil {
				t.Errorf("want error; got no error")
			}

			if accountID != c.want {
				t.Errorf("wrong account id. want: %q, got: %q", c.want, accountID)
			}
		})
	}
}
