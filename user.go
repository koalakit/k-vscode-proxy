package main

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"
)

var gUserCache KVCache[*UserData]

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

// 存储用户数据
func (user *UserData) Save() (err error) {
	savePath := UserDataPath(user.UID)
	user.UpdatedAt = time.Now()
	err = EncodeYamlFile(savePath, user)

	if err != nil {
		LogDebug("save user:", savePath, err)
	} else {
		LogDebug("save user:", savePath)
	}
	return
}

// 加载用户数据
func (user *UserData) Load(openID string) (err error) {
	savePath := UserDataPath(openID)
	err = DecodeYamlFile(savePath, user)

	if err != nil {
		LogDebug("load user:", savePath, err)
	} else {
		LogDebug("load user:", savePath)
	}
	return
}

// 是否存在
func (user *UserData) IsExits(openID string) bool {
	savePath := UserDataPath(openID)
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		return false
	}

	return true
}

func UserDataNew(uid string) (user *UserData) {
	user = new(UserData)
	user.UID = uid
	user.CreatedAt = time.Now()
	user.UserToken = TokenNew()

	return
}

func UserDataPath(uid string) string {
	return path.Join(gAppConfig.UserFolder, fmt.Sprintf("%v.yaml", uid))
}

func UserLoad(uid string) (user *UserData, err error) {
	cacheUser, ok := gUserCache.Get(uid)
	LogDebug("UserLoad:", cacheUser, ok)
	if ok {
		if cacheUser != nil {
			// 命中缓存
			user = cacheUser
			LogDebug("命中缓存: uid:", user.UID)
			return
		}

		gUserCache.Del(uid)
	}

	user = new(UserData)
	if err = user.Load(uid); err != nil {
		return
	}

	gUserCache.Set(uid, user)

	return
}

func UserSave(user *UserData) (err error) {
	if user == nil {
		return
	}

	user.UpdatedAt = time.Now()
	if err = user.Save(); err != nil {
		return
	}

	gUserCache.SetEx(user.UID, user, 7*24*int64(time.Hour))

	return
}

func UserDebug() {
	gUserCache.itemsRW.RLock()
	defer gUserCache.itemsRW.RUnlock()

	LogDebug("UserDebug============================================================")
	for k, v := range gUserCache.items {
		LogDebugf("uid:%s name:%s token:%s", k, v.Value.Name, v.Value.UserToken)
	}
	LogDebug("UserDebug============================================================")
}

func init() {
	gUserCache.Start()
}
