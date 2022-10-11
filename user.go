package main

import (
	"errors"
	"fmt"
	"time"
)

var (
	ErrUserNotFound   = errors.New("user not found")
	ErrUserNil        = errors.New("user is nil")
	ErrUserLoadFailed = errors.New("user load failed")
)

// 用户数据
type UserData struct {
	UID      string `yaml:"uid"`
	FeishuID string `yaml:"feishu-id"`

	Name      string `yaml:"name"`
	AvatarURL string `yaml:"avatar-url"`
	Email     string `yaml:"email"`
	Mobile    string `yaml:"mobile"`

	BackendAddress string `yaml:"backend-address"`

	CreateIP string `yaml:"create-ip"`
	LoginIP  string `yaml:"login-ip"`

	UserToken string `yaml:"user-token"`

	CreatedAt time.Time `yaml:"created-at"`
	UpdatedAt time.Time `yaml:"updated-at"`
}

// 返回用户缓存键名
func UserKey(uid string) string {
	return fmt.Sprintf("user:%v", uid)
}

func UserDataNew(uid string) (user *UserData) {
	user = new(UserData)
	user.UID = uid
	user.CreatedAt = time.Now()
	user.UserToken = TokenNew()

	return
}

func UserLoad(uid string) (user *UserData, err error) {
	var userData UserData
	userData, err = RedisGetJSON[UserData](UserKey(uid))
	if err != nil {
		return
	}

	user = &userData
	return
}

func UserSave(user *UserData) (err error) {
	if user == nil {
		return
	}

	user.UpdatedAt = time.Now()
	err = RedisSetJSON(UserKey(user.UID), user, 0)
	return
}
