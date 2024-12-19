package auth

import (
	"context"
	"time"

	"github.com/auth0/go-jwt-middleware/v2/validator"
)

type Claims struct {
	Username string `json:"username"`
	validator.RegisteredClaims
}

// validator.CustomClaims 接口实现(auth0/go-jwt-middleware)
// 当 JWT 验证通过后会调用此方法
// 当前无需实现
func (c *Claims) Validate(ctx context.Context) error {
	return nil
}

// jwt.Claims 接口实现(golang-jwt/jwt)
func (c *Claims) Valid() error {
	return nil
}

type Token struct {
	Token    string    `json:"token"`
	ExpireAt time.Time `json:"expireAt"`
}
