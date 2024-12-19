package web

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"liewell.fun/alioth/auth"
	"liewell.fun/alioth/core"
	"liewell.fun/alioth/rplace"
	"liewell.fun/alioth/web/middleware"
)

func registryHandler(engine *gin.Engine) {

	engine.GET("/login", auth.Login)
	engine.GET("/register", auth.Register)

	// r/place 请求
	engine.GET("/rplace", rplace.HandleWebSocket)

	// API 请求必须要求验证 JWT
	api := engine.Group("/api", middleware.JWT())
	api.GET("/", func(ctx *gin.Context) {
		username := ctx.MustGet(core.UserNameKey)
		ctx.JSON(http.StatusOK, username)
	})
}

func StartAndWait(ctx context.Context) {

	cfg := &core.GlobalConfig.Http

	// 服务运行模式
	if core.GlobalConfig.Zap.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// 新建实例并配置全局中间件
	r := gin.New()
	r.Use(middleware.Logger(nil), middleware.CORS(), middleware.Recovery())

	// 注册路由
	registryHandler(r)

	// 启动服务
	if len(cfg.ListenTLS) > 0 {
		go func() {
			if err := r.RunTLS(cfg.ListenTLS, cfg.CertFile, cfg.KeyFile); err != nil {
				core.Logger.Fatalf("[web] https server error: %v", err)
			}
		}()
	}
	if len(cfg.Listen) > 0 {
		go func() {
			if err := r.Run(cfg.Listen); err != nil {
				core.Logger.Fatalf("[web] http server error: %v", err)
			}
		}()
	}
	core.Logger.Infof("[web] start success with http[%v], https[%v]", cfg.Listen, cfg.ListenTLS)

	<-ctx.Done()
	core.Logger.Fatalf("[web] server shutdown: %v", ctx.Err())
}
