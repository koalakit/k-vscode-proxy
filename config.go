package main

import (
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

const (
	AppVersion               = "k-vscode-proxy v0.1.0"
	DefaultFeishuRedirectURL = "https://open.feishu.cn/open-apis/authen/v1/index"
	DefaultCookie            = "k-vscode-proxy-token"
)

type AppConfig struct {
	// 配置文件
	ConfigPath string `yaml:"-"`
	// 根目录
	RootFolder string `yaml:"-"`
	// 用户数据目录
	UserFolder string `yaml:"-"`

	// Cookie字段名
	Cookie string `yaml:"cookie"`

	// 飞书配置
	FeishuAppID       string `yaml:"feishu-app-id"`
	FeishuAppSecret   string `yaml:"feishu-app-secret"`
	FeishuAuthenURL   string `yaml:"feishu-authen-url"`
	FeishuRedirectURL string `yaml:"feishu-redirect-url"`
}

func (config *AppConfig) SetRoot(folder string) {
	config.RootFolder = folder
	config.UserFolder = path.Join(config.RootFolder, "user")
	config.ConfigPath = path.Join(config.RootFolder, "app-config.yaml")
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
	gAppConfig.SetRoot(path.Join(os.Getenv("HOME"), ".k-vscode-proxy"))
	gAppConfig.FeishuAuthenURL = DefaultFeishuRedirectURL
	gAppConfig.Cookie = DefaultCookie
}
