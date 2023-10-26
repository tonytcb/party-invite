package cache

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/tonytcb/party-invite/pkg/infrastructure/logger"
	"os"
	"testing"
)

func TestNewInMemoryFilterCustomersCache(t *testing.T) {
	t.Parallel()

	var ctx = context.Background()
	var log = logger.NewLogger(os.Stdout)
	var c = NewInMemoryFilterCustomersCache(log)
	var reader1 = []byte(`{1}`)

	result1, err1 := c.Get(ctx, reader1)
	assert.Nil(t, err1)
	assert.Nil(t, result1)

	err2 := c.Save(ctx, reader1, []byte(`response 1`))
	assert.Nil(t, err2)

	result3, err2 := c.Get(ctx, reader1)
	assert.Nil(t, err2)
	assert.Equal(t, result3, []byte(`response 1`))
}
