package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	kTokenLength = 32
)

const kTokenExpire = 1 * 24 * time.Hour

type TokenData struct {
	UID        int    `yaml:"uid"`
	BackendURL string `yaml:"backend"`
}

func TokenKey(token string) string {
	return fmt.Sprintf("token:%v", token)
}

func TokenSet(token string, uid int, backendURL string) (err error) {
	var data TokenData
	data.UID = uid
	data.BackendURL = backendURL

	err = RedisSetYaml(TokenKey(token), data, kTokenExpire)
	return
}

func TokenGet(token string) (uid int, backendURL string, ok bool) {
	data, err := RedisGetYaml[TokenData](TokenKey(token))
	if err != nil {
		ok = false
		return
	}

	uid = data.UID
	backendURL = data.BackendURL
	ok = true
	return
}

func TokenDel(token string) {
	RedisDel(token)
}

func TokenCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     gAppConfig.Cookie,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(w, cookie)
}

func TokenNew() (token string) {
	rand.Seed(time.Now().UnixNano())
	token = RandomString(kTokenLength)
	return
}
