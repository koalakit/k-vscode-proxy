package main

import (
	"math/rand"
	"net/http"
	"time"
)

const (
	kTokenLength = 32
)

var gTokenCache KVCache[string]

const kTokenExpire = 1 * 24 * int64(time.Hour)

func TokenSet(token string, uid string) {
	gTokenCache.SetEx(token, uid, kTokenExpire)
}

func TokenGet(token string) (uid string, ok bool) {
	uid, ok = gTokenCache.Get(token)
	if !ok {
		return
	}

	// 清理空的缓存
	if len(uid) <= 0 {
		TokenDel(token)

		uid = ""
		ok = false
	}

	return
}

func TokenDel(token string) {
	gTokenCache.Del(token)
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

func TokenDebug() {
	gTokenCache.itemsRW.RLock()
	defer gTokenCache.itemsRW.RUnlock()

	LogDebug("TokenDebug============================================================")
	for k, v := range gTokenCache.items {
		LogDebugf("token: %s:%s", k, v.Value)
	}
	LogDebug("TokenDebug============================================================")
}

func TokenNew() (token string) {
	rand.Seed(time.Now().UnixNano())
	token = RandomString(kTokenLength)
	return
}

func init() {
	gTokenCache.Start()
}
