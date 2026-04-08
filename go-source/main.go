package main

import (
	"fmt"
	"go-source/config"
	"go-source/internal/database"
	"go-source/internal/handler"
	"go-source/router"
	"log"
)

func main() {
	cfg := config.Load()

	// 未安装时跳过数据库连接，直接启动服务等待安装引导
	if handler.IsInstalled() {
		if err := database.Init(cfg); err != nil {
			log.Fatalf("数据库连接失败: %v", err)
		}
		fmt.Println("数据库连接成功")
	} else {
		fmt.Println("系统未安装，请访问 /install 进行安装引导")
	}

	r := router.Setup()
	addr := ":" + cfg.App.Port
	fmt.Printf("服务启动在 http://localhost%s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}
