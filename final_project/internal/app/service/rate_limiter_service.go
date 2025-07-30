package service

import (
	"context"
	"fmt"
	"log"

	"rate_limiter/config"
	"rate_limiter/internal/app/repository"
)

type RateLimiterService struct {
	iPService      IPService
	bucketService  *BucketService
	listRepository repository.ListRepository
}

func NewRateLimiterService(ctx context.Context, config config.AppConfig) (*RateLimiterService, error) {
	repo, err := repository.NewListRepository(config.EnvConfig)
	if err != nil {
		return nil, err
	}

	config.IPWhitelist, _ = repo.GetWhitelist(ctx)
	config.IPBlacklist, _ = repo.GetBlacklist(ctx)
	fmt.Printf("updated Config: %+v\n", config)
	ipService, err := NewIPService(config)
	if err != nil {
		return nil, err
	}

	return &RateLimiterService{
		iPService:      *ipService,
		bucketService:  NewBucketService(config),
		listRepository: *repo,
	}, nil
}

func (r *RateLimiterService) Check(ctx context.Context, login, pass, ip string) bool {
	_ = ctx
	res, err := r.iPService.IsInWhitelist(ip)
	if err != nil {
		return false
	}
	if res {
		log.Printf("IP %s in Whitelist\n", ip)
		return true
	}
	log.Printf("IP %s not in Whitelist\n", ip)
	res, err = r.iPService.IsInBlacklist(ip)
	if err != nil {
		return false
	}
	if res {
		log.Printf("IP %s in Blacklist\n", ip)
		return false
	}
	log.Printf("IP %s not in Blacklist\n", ip)

	return r.bucketService.Check(login, pass, ip)
}

func (r *RateLimiterService) AddToBlacklist(ctx context.Context, networkStr string) (bool, error) {
	res, err := r.iPService.AddToBlacklist(networkStr)
	if err != nil {
		return false, err
	}
	if !res {
		return false, nil
	}

	err = r.listRepository.AddToBlacklist(ctx, networkStr)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *RateLimiterService) RemoveFromBlacklist(ctx context.Context, networkStr string) (bool, error) {
	res, err := r.iPService.RemoveFromBlacklist(networkStr)
	if err != nil {
		return false, err
	}
	if !res {
		return false, nil
	}

	err = r.listRepository.RemoveFromBlacklist(ctx, networkStr)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *RateLimiterService) AddToWhitelist(ctx context.Context, networkStr string) (bool, error) {
	res, err := r.iPService.AddToWhitelist(networkStr)
	if err != nil {
		return false, err
	}
	if !res {
		return false, nil
	}
	log.Printf("Whitelist: %+v", r.iPService.Whitelist)
	err = r.listRepository.AddToWhitelist(ctx, networkStr)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *RateLimiterService) RemoveFromWhitelist(ctx context.Context, networkStr string) (bool, error) {
	res, err := r.iPService.RemoveFromWhitelist(networkStr)
	if err != nil {
		return false, err
	}
	if !res {
		return false, nil
	}

	err = r.listRepository.RemoveFromWhitelist(ctx, networkStr)
	if err != nil {
		return false, err
	}

	return true, nil
}

func (r *RateLimiterService) ResetBucket(ctx context.Context, login, ip string) {
	_ = ctx
	r.bucketService.ResetByLogin(login)
	r.bucketService.ResetByIP(ip)
}
