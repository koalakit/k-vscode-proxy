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

	// 创建工作目录
	// os.MkdirAll(gUserPath, os.ModeDir)

	switch flag.Args()[0] {
	case "serve": // 启动服务
		commandServe(flag.Args()[1:])
	case "user": // 用户设置
		commandUser(flag.Args()[1:])
	case "config":
		commandConfig(flag.Args()[1:])
	default:
		fmt.Println("unknown command")
	}
}

func commandServe(args []string) {
	// 子命令参数解析
	flags := flag.NewFlagSet("serve", flag.ExitOnError)

	debug := flags.Bool("debug", false, "debug mode")
	root := flags.String("root", gAppConfig.RootFolder, "root folder")
	address := flags.String("address", ":8001", "bind address")

	flags.Parse(args)

	// 设置调试模式
	if *debug {
		fmt.Println(AppVersion, "(debug mode)")
		LogMode(true)
	}

	var err error

	gAppConfig.SetRoot(*root)
	if err = DecodeYamlFile(gAppConfig.ConfigPath, &gAppConfig); err != nil {
		fmt.Println(err)
		return
	}

	// 初始化数据库
	if err = RedisInit(gAppConfig.RedisDB); err != nil {
		log.Fatal("[Redis]", err)
		return
	}

	var app = new(ProxyApp)
	var server = http.Server{
		Addr:    *address,
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

func commandUser(args []string) {
	// 子命令参数解析
	flags := flag.NewFlagSet("user", flag.ExitOnError)

	uid := flags.String("uid", "", "user id of user data")
	fieldName := flags.String("field", "", "field name of user data")
	fieldValue := flags.String("field-value", "", "field value of user data")

	flags.Parse(args)

	var err error

	// gAppConfig.SetRoot(*root)
	if err = DecodeYamlFile(gAppConfig.ConfigPath, &gAppConfig); err != nil {
		fmt.Println(err)
		return
	}

	// 初始化数据库
	if err = RedisInit(gAppConfig.RedisDB); err != nil {
		log.Fatal("[Redis]", err)
		return
	}

	var userData map[string]any

	if len(*uid) <= 0 {
		fmt.Println("缺少用户ID --uid")
		return
	}

	if len(*fieldName) <= 0 {
		fmt.Println("缺少字段名 --field")
		return
	}

	// 设置字段值
	if len(*fieldValue) > 0 {
		err = RedisGetJSONEx("user:"+*uid, &userData)
		if err != nil {
			log.Fatal(err)
			return
		}

		userData[*fieldName] = *fieldValue
		err = RedisSetJSON("user:"+*uid, userData, 0)
		if err != nil {
			log.Fatal(err)
			return
		}
		return
	}

	// 读取字段值
	{
		if err = RedisGetJSONEx("user:"+*uid, &userData); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%s: %v\n", *fieldName, userData[*fieldName])
	}
}

func commandConfig(args []string) {
	// 子命令参数解析
	flags := flag.NewFlagSet("user", flag.ExitOnError)

	init := flags.Bool("init", false, "initialize configuration and root folder")
	root := flags.String("root", gAppConfig.RootFolder, "root folder")
	fieldName := flags.String("field", "", "field name of user data")
	fieldValue := flags.String("field-value", "", "field value of user data")

	flags.Parse(args)

	var err error

	gAppConfig.SetRoot(*root)

	if *init {
		// 创建根目录
		os.Mkdir(gAppConfig.RootFolder, os.ModeDir)

		// 保存默认配置
		if err = EncodeYamlFile(gAppConfig.ConfigPath, &gAppConfig); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(gAppConfig.ConfigPath)
		return
	}

	if len(*fieldName) <= 0 {
		fmt.Println("缺少字段名 --field")
		return
	}

	var configData map[string]any
	// 设置字段值
	if len(*fieldValue) > 0 {
		if err = DecodeYamlFile(gAppConfig.ConfigPath, &configData); err != nil {
			fmt.Println(err)
			return
		}

		configData[*fieldName] = *fieldValue
		if err = EncodeYamlFile(gAppConfig.ConfigPath, configData); err != nil {
			fmt.Println(err)
		}
		return
	}

	// 读取字段值
	{
		if err = DecodeYamlFile(gAppConfig.ConfigPath, &configData); err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%s: %v\n", *fieldName, configData[*fieldName])
	}
}
