package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"liewell.fun/alioth/core"
	"liewell.fun/alioth/web"
)

func main() {

	// 全局 context 控制
	ctx, cancel := context.WithCancel(context.Background())
	go WaitTerm(cancel)

	// 读取配置文件并解析
	c := flag.String("c", "config.yaml", "config file path")
	flag.Parse()
	core.LoadYamlConfig(*c)

	// 初始化日志模块,必须是第一个被初始化的模块
	core.InitZap()

	// 初始化数据库
	core.InitMysql()

	// 启动 web 服务
	web.StartAndWait(ctx)
}

func WaitTerm(cancel context.CancelFunc) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(quit)
	<-quit
	cancel()
}
