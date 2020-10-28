package cache_test

import (
	"Nani/internal/app/cache"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

type mockCache struct {
	cache map[string]interface{}
}

func (m *mockCache) Set(key string, value interface{}) {
	m.cache[key] = value
}

func (m *mockCache) GetV(key string) (interface{}, error) {
	v, ok := m.cache[key]
	if !ok {
		return nil, fmt.Errorf("error key")
	}

	return v, nil
}

func CreateCache() *mockCache {
	return &mockCache{
		cache: make(map[string]interface{}),
	}
}

func TestSet_ShouldSetNewKeywordToCache_NoError(t *testing.T) {
	kc := cache.NewKeyCache(CreateCache())
	assert.NoError(t, kc.Set("key"))
}

func TestNext_ShouldCreateNewValueForNextKey_NoError(t *testing.T) {
	kc := cache.NewKeyCache(CreateCache())
	assert.NoError(t, kc.Set("key"))
	str, err := kc.Next()
	assert.NoError(t, err)
	assert.Equal(t, "key", str)
}

func TestNext_ShouldCreateNewSomeValueForKeys_NoError(t *testing.T) {
	kc := cache.NewKeyCache(CreateCache())
	assert.NoError(t, kc.Set("key"))
	assert.NoError(t, kc.Set("key1"))
	assert.NoError(t, kc.Set("key2"))
	str, err := kc.Next()
	assert.NoError(t, err)
	assert.Equal(t, "key", str)
	str, err = kc.Next()
	assert.NoError(t, err)
	assert.Equal(t, "key1", str)
	str, err = kc.Next()
	assert.NoError(t, err)
	assert.Equal(t, "key2", str)
}

func TestNext_ShouldReturnEmptyStringCozCacheIsEmpty_Error(t *testing.T) {
	kc := cache.NewKeyCache(CreateCache())
	str, err := kc.Next()
	assert.Error(t, err)
	assert.Empty(t, str)
	str, err = kc.Next()
	assert.Error(t, err)
	assert.Empty(t, str)
	str, err = kc.Next()
	assert.Error(t, err)
	assert.Empty(t, str)
}

func TestNext_ShouldReturnErrorIfIndexOutOfRange_Error(t *testing.T) {
	kc := cache.NewKeyCache(CreateCache())
	assert.NoError(t, kc.Set("kry1"))
	str, err := kc.Next()
	assert.NoError(t, err)
	assert.NotEmpty(t, str)
	str, err = kc.Next()
	assert.Error(t, err)
	assert.Empty(t, str)
}

func TestRollback_ShouldRollbackToPrevPosition_NoErrors(t *testing.T) {
	kc := cache.NewKeyCache(CreateCache())
	assert.NoError(t, kc.Set("kry1"))
	assert.NoError(t, kc.Set("kry2"))
	assert.NoError(t, kc.Set("kry3"))
	str, err := kc.Next()
	assert.NoError(t, err)
	assert.Equal(t,"kry1", str)
	str, err = kc.Next()
	assert.NoError(t, err)
	assert.Equal(t,"kry2", str)
	assert.NoError(t, kc.Rollback())
	str, err = kc.Next()
	assert.NoError(t, err)
	assert.Equal(t,"kry2", str)
}

func TestRollback_ShouldReturnNoErrorIfCacheIsEmpty_NoError(t *testing.T) {
	kc := cache.NewKeyCache(CreateCache())
	assert.NoError(t, kc.Set("key1"))
	assert.NoError(t, kc.Rollback())
}




