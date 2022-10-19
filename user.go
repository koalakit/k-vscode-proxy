package main

import (
	"errors"
	"time"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserNil        = errors.New("user is nil")
	ErrUserLoadFailed = errors.New("user load failed")
)

// 用户数据
type UserData struct {
	ID       int    `gorm:"primaryKey" yaml:"id"`
	FeishuID string `gorm:"size:64;index" yaml:"feishu-id"`

	Account   string `gorm:"size:16;index" yaml:"account"`
	Name      string `gorm:"size:16" yaml:"name"`
	AvatarURL string `gorm:"size:256" yaml:"avatar-url"`
	Email     string `gorm:"size:32" yaml:"email"`
	Mobile    string `gorm:"size:16" yaml:"mobile"`

	BackendAddress string `gorm:"size:256" yaml:"backend-address"`

	CreateIP string `gorm:"size:16;column:create_ip" yaml:"create-ip"`
	LastIP   string `gorm:"size:16;column:last_ip" yaml:"login-ip"`

	UserToken string `gorm:"size:32" yaml:"user-token"`

	CreatedAt time.Time `yaml:"created-at"`
	UpdatedAt time.Time `yaml:"updated-at"`
}

func (user *UserData) TableName() string {
	return "user"
}

// UserDataCreateFromFeishu 通过飞书OpenID创建用户
func UserDataCreateFromFeishu(feishuID string) (user *UserData, err error) {
	user = new(UserData)
	user.FeishuID = feishuID

	db := mysqldb.Create(user)
	err = db.Error
	return
}

// UserLoadFromFeishu 通过飞书OpenID加载用户
func UserLoadFromFeishu(feishuID string) (user *UserData, err error) {
	user = new(UserData)

	db := mysqldb.Raw("SELECT * FROM user WHERE feishu_id=? LIMIT 1", feishuID).Scan(user)
	if err = db.Error; err != nil {
		user = nil
		return
	}

	return
}

// UserLoad 通过账号加载用户
func UserLoad(account string) (user *UserData, err error) {
	user = new(UserData)

	db := mysqldb.Raw("SELECT * FROM user WHERE account=? LIMIT 1", account).Scan(user)
	if err = db.Error; err != nil {
		user = nil
		return
	}

	return
}

// UserSave 保存用户数据
func UserSave(user *UserData) (err error) {
	db := mysqldb.Save(user)
	err = db.Error
	return
}
