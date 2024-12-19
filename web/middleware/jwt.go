package middleware

import (
	"context"
	"net/http"
	"time"

	"liewell.fun/alioth/auth"
	"liewell.fun/alioth/core"

	jwtmiddleware "github.com/auth0/go-jwt-middleware/v2"
	"github.com/auth0/go-jwt-middleware/v2/validator"
	"github.com/gin-gonic/gin"
)

func JWT() gin.HandlerFunc {

	cfg := core.GlobalConfig.JWT

	// 定义验证器,验证器配置失败时退出应用
	jwtValidator, err := validator.New(
		func(ctx context.Context) (interface{}, error) {
			return []byte(cfg.Secret), nil
		},
		validator.HS256,
		cfg.Issuer,
		cfg.Audience,
		validator.WithCustomClaims(func() validator.CustomClaims {
			return &auth.Claims{}
		}), // 自定义验证对象
		validator.WithAllowedClockSkew(5*time.Second), // 时钟偏差
	)
	if err != nil {
		core.Logger.Panic("[JWT] set up validator failed: %v", err)
	}
	middleware := jwtmiddleware.New(
		jwtValidator.ValidateToken,
		jwtmiddleware.WithErrorHandler(func(w http.ResponseWriter, r *http.Request, err error) {
			// 自定义错误仅记录原因,不将错误返回给调用方
			core.Logger.Info("[JWT] validate error: %v", err)
		}))

	return func(ctx *gin.Context) {
		encounteredError := true
		var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
			encounteredError = false
			ctx.Request = r

			// 当 JWT 与自定义 Claims 检验均通过时, 将用户信息写入上下文
			core.Logger.Info("[JWT] validate success, set username to context")
			claims, _ := ctx.Request.Context().Value(jwtmiddleware.ContextKey{}).(*validator.ValidatedClaims)
			customClaims, _ := claims.CustomClaims.(*auth.Claims)
			ctx.Set(core.UserNameKey, customClaims.Username)

			// 交给后续 handler 继续处理
			ctx.Next()
		}

		// 校验 JWT
		middleware.CheckJWT(handler).ServeHTTP(ctx.Writer, ctx.Request)

		// 如果校验失败,则返回错误并中断处理
		if encounteredError {
			ctx.AbortWithStatusJSON(
				http.StatusUnauthorized,
				core.UnauthorizedError(auth.JWTErrorCode, auth.JWTErrorMessage),
			)
		}
	}
}
