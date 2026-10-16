package web

import (
	"context"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"liewell.fun/alioth/auth"
	"liewell.fun/alioth/core"
	"liewell.fun/alioth/rplace"
	"liewell.fun/alioth/static"
	"liewell.fun/alioth/web/middleware"
)

func registryHandler(engine *gin.Engine) {

	// API/WS 路由先注册（NoRoute 兜底在后，避免吞掉 /rplace WebSocket 升级）
	engine.GET("/login", auth.Login)
	engine.GET("/register", auth.Register)
	engine.GET("/rplace", rplace.HandleWebSocket)

	// API 请求必须要求验证 JWT
	api := engine.Group("/api", middleware.JWT())
	api.GET("/", func(ctx *gin.Context) {
		username := ctx.MustGet(core.UserNameKey)
		ctx.JSON(http.StatusOK, username)
	})

	// 静态资源（前端页面，go:embed 嵌入到 alioth/static 包）
	// 根级文件（index/about/favicon/robots）走 embed 根 FS
	// css/js 各自 sub 到对应目录，因为 gin StaticFS 会剥离 URL 前缀后直接 fs.Open
	sub := static.Sub()
	hfs := http.FS(sub)
	cssFS, _ := fs.Sub(static.FS, "css")
	jsFS, _ := fs.Sub(static.FS, "js")
	serveFile := func(name string) gin.HandlerFunc {
		return func(c *gin.Context) {
			f, err := hfs.Open(name)
			if err != nil {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			defer f.Close()
			c.Status(http.StatusOK)
			c.Header("Content-Type", "text/html; charset=utf-8")
			_, _ = io.Copy(c.Writer, f)
		}
	}
	engine.GET("/", serveFile("index.html"))
	engine.GET("/about", serveFile("about.html"))
	engine.StaticFS("/css", http.FS(cssFS))
	engine.StaticFS("/js", http.FS(jsFS))
	engine.StaticFileFS("/favicon.ico", "favicon.ico", hfs)
	engine.StaticFileFS("/robots.txt", "robots.txt", hfs)

	// SPA 兜底：未匹配的路径返回 index.html
	engine.NoRoute(func(c *gin.Context) {
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(http.StatusNotFound, gin.H{"code": "404", "message": http.StatusText(http.StatusNotFound)})
			return
		}
		serveFile("index.html")(c)
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
	r.RedirectTrailingSlash = false // 关闭尾斜杠 301，让未匹配请求直接走 NoRoute（SPA 兜底）
	r.RedirectFixedPath = false     // 关闭路径修正 301
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
