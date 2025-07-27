package service

import (
	"sync"
	"time"

	"rate_limiter/config"
	"rate_limiter/internal/app/bucket"
	"rate_limiter/internal/app/cache"
)

type BucketService struct {
	K int
	M int
	N int

	BucketStorage bucket.Storage
	lock          sync.Mutex
}

func NewBucketService(config config.AppConfig) *BucketService {
	bucketStorage := bucket.NewMemoryStorage()
	bucketStorage.Init()
	return &BucketService{
		K: config.K,
		M: config.M,
		N: config.N,

		BucketStorage: bucketStorage,
	}
}

func (r *BucketService) Check(login, pass, ip string) bool {
	res := r.checkBucket(bucket.LoginType, login)
	if !res {
		return false
	}

	res = r.checkBucket(bucket.PasswordType, pass)
	if !res {
		return false
	}

	return r.checkBucket(bucket.IPType, ip)
}

func (r *BucketService) checkBucket(bucketType, keyStr string) bool {
	key := cache.Key(keyStr)
	curBucket := r.getOrCreate(bucketType, key)
	res := (*curBucket).Add()
	_ = r.BucketStorage.UpdateBucket(bucketType, key, curBucket)
	return res
}

func (r *BucketService) getOrCreate(cacheType string, key cache.Key) *bucket.LeakyBucket {
	r.lock.Lock()
	defer r.lock.Unlock()

	existingBucket := r.BucketStorage.GetBucketByKey(cacheType, key)
	if existingBucket != nil {
		// log.Printf("existing bucket: %+v\n", *existingBucket)
		return existingBucket
	}

	capacity := 0
	switch cacheType {
	case bucket.LoginType:
		capacity = r.K
	case bucket.PasswordType:
		capacity = r.M
	case bucket.IPType:
		capacity = r.N
	}
	newBucket := bucket.NewBucket(capacity, time.Minute)
	_ = r.BucketStorage.UpdateBucket(cacheType, key, &newBucket)

	return &newBucket
}

func (r *BucketService) ResetByLogin(login string) {
	loginBucket := r.BucketStorage.GetBucketByKey(bucket.LoginType, cache.Key(login))
	if loginBucket == nil {
		return
	}
	(*loginBucket).Reset()
}

func (r *BucketService) ResetByIP(ip string) {
	ipBucket := r.BucketStorage.GetBucketByKey(bucket.IPType, cache.Key(ip))
	if ipBucket == nil {
		return
	}
	(*ipBucket).Reset()
}
