package entity

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/golang-jwt/jwt/v5"
	"github.com/shiqinfeng1/goframe-ddd/pkg/clock"
)

// t.go (值对象)
type Token struct {
	UserID    string `json:"user_id"`
	RefreshID string `json:"refresh_id"`
}

func NewToken() *Token {
	return &Token{}
}

// 生成JWT Token
func generateJWT(secret string, claims *jwt.RegisteredClaims) (string, error) {
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(secret))
}

// 生成Token对
func (s *Token) GenerateTokenPair(ctx context.Context, userId string) (accessToken, refreshToken string, err error) {
	now := clock.Now()
	s.UserID = userId
	accessSecret := g.Cfg().MustGet(ctx, "jwt.accessSecret").String()

	// 生成Access Token
	accessToken, err = generateJWT(
		accessSecret,
		&jwt.RegisteredClaims{
			Subject:   gconv.String(s.UserID),
			Issuer:    "go-mgrid",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(g.Cfg().MustGet(ctx, "jwt.accessExpire").Duration())),
		},
	)
	if err != nil {
		return "", "", err
	}

	// 生成Refresh Token (带唯一ID)
	s.RefreshID = guid.S()
	refreshToken, err = generateJWT(
		accessSecret,
		&jwt.RegisteredClaims{
			Subject:   gconv.String(s.UserID),
			ID:        s.RefreshID,
			Issuer:    "go-mgrid",
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(g.Cfg().MustGet(ctx, "jwt.refreshExpire").Duration())),
		},
	)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, err
}

// 解析JWT Token
func (s *Token) ParseJWT(tokenString, secret string) error {
	t, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !t.Valid {
		return gerror.New("invalid jwt token")
	}

	claims := t.Claims.(jwt.RegisteredClaims)
	s.UserID = claims.Subject
	s.RefreshID = claims.ID
	return nil
}
