package main

import (
	"sync"
	"time"
)

type KVCacheItem[ValueT any] struct {
	Value  ValueT
	Expire int64
}

func (item *KVCacheItem[ValueT]) IsExpire(now int64) bool {
	if item.Expire == 0 {
		return false
	}

	if now == 0 {
		now = time.Now().Unix()
	}

	if now > item.Expire {
		return true
	}

	return false
}

type KVCache[ValueT any] struct {
	items   map[string]*KVCacheItem[ValueT]
	itemsRW sync.RWMutex
}

func (cache *KVCache[ValueT]) Start() {
	cache.items = make(map[string]*KVCacheItem[ValueT])

	go func() {
		ticker := time.NewTicker(10 * time.Minute)
		for {
			<-ticker.C
			cache.GC()
			Log("cache gc.")
		}
	}()
}

func (cache *KVCache[ValueT]) Get(k string) (v ValueT, ok bool) {
	cache.itemsRW.RLock()
	defer cache.itemsRW.RUnlock()

	item, ok := cache.items[k]
	if !ok {
		return
	}

	// 超时判断
	if item.IsExpire(0) {
		ok = false
		return
	}

	v = item.Value

	return
}

func (cache *KVCache[ValueT]) Set(k string, v ValueT) {
	cache.itemsRW.Lock()
	defer cache.itemsRW.Unlock()

	item := new(KVCacheItem[ValueT])
	item.Value = v
	item.Expire = 0

	cache.items[k] = item
}

func (cache *KVCache[ValueT]) SetEx(k string, v ValueT, Expire int64) {
	cache.itemsRW.Lock()
	defer cache.itemsRW.Unlock()

	item := new(KVCacheItem[ValueT])
	item.Value = v
	item.Expire = time.Now().Add(7 * 24 * time.Hour).Unix()

	cache.items[k] = item
}

func (cache *KVCache[ValueT]) Del(k string) {
	cache.itemsRW.Lock()
	defer cache.itemsRW.Unlock()

	delete(cache.items, k)
}

func (cache *KVCache[ValueT]) GC() {
	cache.itemsRW.Lock()
	defer cache.itemsRW.Unlock()

	now := time.Now().Unix()
	expireItems := make([]string, 0, len(cache.items))

	for k, v := range cache.items {
		if v.IsExpire(now) {
			expireItems = append(expireItems, k)
		}
	}

	for _, v := range expireItems {
		delete(cache.items, v)
	}
}
