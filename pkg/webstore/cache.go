package webstore

import "github.com/sir-wiggles/arc/pkg/webstore/redis"

type Cache struct {
	service *redis.Service
}

type CacheService interface {
	Set(string, string) error
	Get(string) (string, error)
	Delete(string) error
}
