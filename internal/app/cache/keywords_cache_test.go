package cache_test

import (
	"Nani/internal/app/cache"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

type mockCache struct {
	cache map[string]interface{}
	mutex sync.Mutex
}

func (m *mockCache) Set(key string, value interface{}) {
	m.mutex.Lock()
	m.cache[key] = value
	m.mutex.Unlock()
}

func (m *mockCache) GetV(key string) (interface{}, error) {
	//m.mutex.Lock()
	v, ok := m.cache[key]
	if !ok {
		return nil, fmt.Errorf("error key")
	}
	//m.mutex.Unlock()
	return v, nil
}

func (m *mockCache) Dump() {}

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

func TestNext_ShouldReturnErrorIfCacheIsEmpty_Error(t *testing.T) {
	kc := cache.NewKeyCache(CreateCache())
	str, err := kc.Next()
	assert.Error(t, err)
	assert.Empty(t, str)
	assert.Equal(t, "keywords cache is empty", err.Error())

	err = kc.Set("key")
	assert.NoError(t, err)

	str, err = kc.Next()
	assert.NoError(t, err)
	assert.NotEmpty(t, str)
	assert.Equal(t, "key", str)

	str, err = kc.Next()
	assert.Error(t, err)
	assert.Empty(t, str)
	assert.Equal(t, "keywords are out of range", err.Error())
}

func TestRollback_ShouldRollbackToPrevPosition_NoErrors(t *testing.T) {
	kc := cache.NewKeyCache(CreateCache())
	assert.NoError(t, kc.Set("kry1"))
	assert.NoError(t, kc.Set("kry2"))
	assert.NoError(t, kc.Set("kry3"))
	str, err := kc.Next()
	assert.NoError(t, err)
	assert.Equal(t, "kry1", str)
	str, err = kc.Next()
	assert.NoError(t, err)
	assert.Equal(t, "kry2", str)
	assert.NoError(t, kc.Rollback())
	str, err = kc.Next()
	assert.NoError(t, err)
	assert.Equal(t, "kry2", str)
}

func TestRollback_ShouldReturnNoErrorIfCacheIsEmpty_NoError(t *testing.T) {
	kc := cache.NewKeyCache(CreateCache())
	assert.NoError(t, kc.Set("key1"))
	assert.NoError(t, kc.Rollback())
}

func TestDistinct_ShouldRemoveDuplicatesFromCache_NoError(t *testing.T) {
	kc := cache.NewKeyCache(CreateCache())

	kc.Set("key1")
	kc.Set("key1")
	kc.Set("key1")
	kc.Set("key2")
	kc.Set("key2")
	kc.Set("key3")
	kc.Set("key3")
	kc.Set("key4")

	assert.NoError(t, kc.Distinct())
	dubls := make(map[string]int)
	for {
		key, err := kc.Next()
		if err != nil {
			break
		}
		dubls[key]+=1
	}

	assert.Equal(t, 4, len(dubls))

	for _, v := range dubls {
		assert.Equal(t, 1, v)
	}
}

func TestDistinct_ShouldReturnErrorCozKeyNotFound_Error(t *testing.T) {
	kc := cache.NewKeyCache(CreateCache())

	assert.Error(t, kc.Distinct())
}

