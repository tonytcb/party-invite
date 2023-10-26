package cache

import (
	"context"
	"crypto/md5"
	"fmt"
	"sync"

	"github.com/pkg/errors"

	"github.com/tonytcb/party-invite/pkg/infrastructure/logger"
)

type InMemoryFilterCustomersCache struct {
	data sync.Map
	log  logger.Logger
}

func NewInMemoryFilterCustomersCache(log logger.Logger) *InMemoryFilterCustomersCache {
	return &InMemoryFilterCustomersCache{
		data: sync.Map{},
		log:  log,
	}
}

func (f *InMemoryFilterCustomersCache) Get(ctx context.Context, fileContents []byte) ([]byte, error) {
	log := f.log.FromContext(ctx)

	key := fileMd5(fileContents)

	content, ok := f.data.Load(key)
	if !ok {
		log.Infof("Cache miss, key=%s", key)

		return nil, nil
	}

	if v, ok := content.([]byte); ok {
		log.Infof("Cache hit, key=%s", key)

		return v, nil
	}

	return nil, errors.New("error to load content on cache")
}

func (f *InMemoryFilterCustomersCache) Save(ctx context.Context, fileContents []byte, response []byte) error {
	key := fileMd5(fileContents)

	f.data.Store(key, response)

	f.log.FromContext(ctx).Infof("Cache updated, key=%s", key)

	return nil
}

func fileMd5(content []byte) string {
	return fmt.Sprintf("%x", md5.Sum(content))
}
