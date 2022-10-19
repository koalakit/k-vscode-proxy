package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var (
	flagVersion = flag.Bool("version", false, "应用版本信息")
)

func main() {
	debug := flag.Bool("debug", false, "debug mode")
	configPath := flag.String("config", "", "app config")
	createConfig := flag.Bool("create-config", false, "save config")

	flag.Parse()

	// 打印版本
	if *flagVersion {
		fmt.Println(AppVersion)
		return
	}

	// 默认打印帮助
	if len(os.Args) < 2 {
		flag.Usage()
		return
	}

	// 调试模式
	if *debug {
		fmt.Println(AppVersion, "(debug mode)")
		LogMode(true)
	}

	var err error

	// 加载配置文件
	if len(*configPath) > 0 {
		if err = DecodeYamlFile(*configPath, &gAppConfig); err != nil {
			fmt.Println(err)
			return
		}
	}

	gAppConfig.LoadEnv()

	if *createConfig {
		if err = EncodeYamlFile(*configPath, gAppConfig); err != nil {
			log.Fatal(err)
			return
		}
	}

	Log(gAppConfig)

	// 初始化Redis
	if err = RedisInit(gAppConfig.RedisURL); err != nil {
		log.Fatal("[Redis]", err)
		return
	}

	// 初始化MySQL
	if err = MySQLInit(gAppConfig.MySQLURL, *debug); err != nil {
		log.Fatal("[MySQL]", err)
		return
	}

	var app = new(ProxyApp)
	var server = http.Server{
		Addr:    gAppConfig.Addr,
		Handler: app,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			Log(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		Log("server shutdown failed:", err)
		return
	}

	Log("server shutdown.")
}
