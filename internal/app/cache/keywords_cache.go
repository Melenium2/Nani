package cache

import (
	"errors"
	"fmt"
)

// KeyStorage interface who manages the instance of KeywordCache
type KeyStorage interface {
	Set(key string) error
	Next() (string, error)
	Rollback() error
}

// KeywordCache manages cache for storing key data
type KeywordsCache struct {
	cache   Storage
	key     string
	next    string
	isEmpty bool
}

// Set new key to the cache
func (kc *KeywordsCache) Set(key string) error {
	kc.isEmpty = false
	keysIn, err := kc.cache.GetV(kc.key)
	keys, ok := keysIn.([]Keyword)
	if err != nil || !ok {
		kc.cache.Set(kc.key, []Keyword{{Pos: len(keys) + 1, Key: key}})
		return nil
	}

	kc.cache.Set(kc.key, append(keys, Keyword{Pos: len(keys) + 1, Key: key}))
	return nil
}

// Get Next key from cache keywords slice
func (kc *KeywordsCache) Next() (string, error) {
	if kc.isEmpty {
		return "", fmt.Errorf("keywords cache is empty")
	}

	keyIn, err := kc.cache.GetV(kc.next)
	if err != nil {
		kc.cache.Set(kc.next, 0)
		k, err := kc.get(0)
		if err != nil {
			return "", err
		}
		return k.Key, nil
	}

	key := getKey(keyIn) + 1

	keyword, err := kc.get(key)
	if err != nil {
		kc.isEmpty = true
		return "", errors.New("keywords are out of range")
	}

	kc.cache.Set(kc.next, key)

	return keyword.Key, nil
}

// Rollback key index to the -1
func (kc *KeywordsCache) Rollback() error {
	keyIn, err := kc.cache.GetV(kc.next)
	if err != nil {
		return nil
	}
	key := getKey(keyIn)

	if key == 0 {
		return nil
	}

	kc.cache.Set(kc.next, key-1)
	return nil
}

// Get item with index i
func (kc *KeywordsCache) get(i int) (Keyword, error) {
	keyIn, err := kc.cache.GetV(kc.key)
	if err != nil {
		return Keyword{}, errors.New("keywords cache is empty")
	}
	key, ok := keyIn.([]Keyword)
	if !ok {
		return Keyword{}, nil
	}

	if i >= len(key) {
		return Keyword{}, errors.New("wrong index")
	}

	return key[i], nil
}

// getKey cast interface to int
func getKey(k interface{}) int {
	var num int
	switch n := k.(type) {
	case float64:
		num = int(n)
	case int:
		num = n
	default:
	}

	return num
}

// Create new instance of keywordsCache
func NewKeyCache(cache Storage) *KeywordsCache {
	c := &KeywordsCache{
		cache: cache,
		key:   "_keys",
		next:  "_keys_next",
		isEmpty: true,
	}
	return c
}
