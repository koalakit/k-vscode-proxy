package main

import (
	"os"

	"gopkg.in/yaml.v3"
)

const (
	AppVersion = "k-vscode-proxy v0.1.0"
)

type AppConfig struct {
	Debug bool   `yaml:"-"`
	Addr  string `yaml:"addr"`
	// Cookie字段名
	Cookie string `yaml:"cookie"`

	// MySQL
	MySQLURL string `yaml:"mysql-url"`
	// 数据库
	RedisURL string `yaml:"redis-url"`

	// 飞书配置
	FeishuAppID       string `yaml:"feishu-app-id"`
	FeishuAppSecret   string `yaml:"feishu-app-secret"`
	FeishuAuthenURL   string `yaml:"feishu-authen-url"`
	FeishuRedirectURL string `yaml:"feishu-redirect-url"`
}

func (config *AppConfig) LoadEnv() {
	if v := os.Getenv("KVP_ADDR"); len(v) > 0 {
		config.Addr = v
	}

	if v := os.Getenv("KVP_MYSQL_URL"); len(v) > 0 {
		config.MySQLURL = v
	}
	if v := os.Getenv("KVP_REDIS_URL"); len(v) > 0 {
		config.RedisURL = v
	}

	if v := os.Getenv("KVP_FEISHU_APP_ID"); len(v) > 0 {
		config.FeishuAppID = v
	}

	if v := os.Getenv("KVP_FEISHU_APP_SECRET"); len(v) > 0 {
		config.FeishuAppSecret = v
	}

	if v := os.Getenv("KVP_FEISHU_AUTHEN_URL"); len(v) > 0 {
		config.FeishuAuthenURL = v
	}

	if v := os.Getenv("KVP_FEISHU_REDIRECT_URL"); len(v) > 0 {
		config.FeishuRedirectURL = v
	}
}

var gAppConfig AppConfig

// DecodeYamlFile 加载yaml文件
func DecodeYamlFile(name string, v interface{}) error {
	bs, err := os.ReadFile(name)
	if err != nil {
		return err
	}

	if err = yaml.Unmarshal(bs, v); err != nil {
		return err
	}

	return nil
}

// EncodeYamlFile 保存yaml文件
func EncodeYamlFile(name string, v interface{}) error {
	var err error
	var data []byte
	if data, err = yaml.Marshal(v); err != nil {
		return err
	}

	if err = os.WriteFile(name, data, os.ModeAppend); err != nil {
		return err
	}

	return nil
}

func init() {
	gAppConfig.FeishuAuthenURL = "https://open.feishu.cn/open-apis/authen/v1/index"
	gAppConfig.Cookie = "k-vscode-proxy-token"
	gAppConfig.Addr = ":8001"
}
