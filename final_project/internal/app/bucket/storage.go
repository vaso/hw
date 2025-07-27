package bucket

import (
	"rate_limiter/internal/app/cache"
)

const (
	LoginType    = "login"
	PasswordType = "password"
	IPType       = "ip"
)

type Storage interface {
	Init()
	GetBucketByKey(cacheType string, key cache.Key) *LeakyBucket
	UpdateBucket(cacheType string, key cache.Key, leakyBucket *LeakyBucket) bool
}

type MemoryStorage struct {
	loginBucketCache cache.Cache
	passBucketCache  cache.Cache
	ipBucketCache    cache.Cache
}

func NewMemoryStorage() Storage {
	return &MemoryStorage{}
}

func (s *MemoryStorage) Init() {
	s.loginBucketCache = cache.NewCache(10000)
	s.passBucketCache = cache.NewCache(10000)
	s.ipBucketCache = cache.NewCache(10000)
}

func (s *MemoryStorage) GetBucketByKey(cacheType string, key cache.Key) *LeakyBucket {
	cacheByType := s.getCache(cacheType)
	bucketVal, ok := cacheByType.Get(key)
	if !ok {
		return nil
	}

	bucket := bucketVal.(LeakyBucket)
	return &bucket
}

func (s *MemoryStorage) UpdateBucket(cacheType string, key cache.Key, leakyBucket *LeakyBucket) bool {
	cacheByType := s.getCache(cacheType)
	res := cacheByType.Set(key, *leakyBucket)
	// log.Printf("updated Cache for type %s: %+v\n", cacheType, cacheByType)

	return res
}

func (s *MemoryStorage) ClearCacheByType(cacheType string) {
	cacheByType := s.getCache(cacheType)
	cacheByType.Clear()
}

func (s *MemoryStorage) getCache(cacheType string) cache.Cache {
	var cacheByType cache.Cache
	switch cacheType {
	case LoginType:
		cacheByType = s.loginBucketCache
	case PasswordType:
		cacheByType = s.passBucketCache
	case IPType:
		cacheByType = s.passBucketCache
	}

	return cacheByType
}
