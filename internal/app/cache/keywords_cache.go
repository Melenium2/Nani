package cache

import (
	"encoding/json"
	"errors"
	"fmt"
)

// KeyStorage interface who manages the instance of KeywordCache
type KeyStorage interface {
	Set(key string) error
	Next() (string, error)
	Rollback() error
	Distinct() error
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

	keys := keywords(keysIn)

	if err != nil {
		kc.cache.Set(kc.key, []Keyword{{Pos: len(keys) + 1, Key: key}})
		return nil
	}

	kc.cache.Set(kc.key, append(keys, Keyword{Pos: len(keys) + 1, Key: key}))
	return nil
}

func (kc *KeywordsCache) Distinct() error {
	if kc.isEmpty {
		return fmt.Errorf("keywords cahce is empty")
	}

	keysIn, err := kc.cache.GetV(kc.key)
	if err != nil {
		return err
	}

	keys := keywords(keysIn)
	distinct := make(map[string]struct{})

	for _, v := range keys {
		distinct[v.Key] = struct{}{}
	}

	result := make([]Keyword, len(distinct))
	i := 0
	for k, _ := range distinct {
		result[i] = Keyword{Pos: i+1, Key: k}
		i++
	}

	kc.cache.Set(kc.key, result)

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
		k, err := kc.item(0)
		if err != nil {
			return "", err
		}
		return k.Key, nil
	}
	key := key(keyIn) + 1

	keyword, err := kc.item(key)
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
	key := key(keyIn)

	if key == 0 {
		return nil
	}

	kc.cache.Set(kc.next, key-1)
	return nil
}

// Get item with index i
func (kc *KeywordsCache) item(i int) (Keyword, error) {
	keyIn, err := kc.cache.GetV(kc.key)
	if err != nil {
		return Keyword{}, errors.New("keywords cache is empty")
	}
	keys := keywords(keyIn)

	if i >= len(keys) {
		return Keyword{}, errors.New("wrong index")
	}

	return keys[i], nil
}

// key cast interface to int
func key(k interface{}) int {
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

func keywords(keyIn interface{}) []Keyword {
	if keyIn == nil {
		return []Keyword{}
	}
	keysJson, _ := json.Marshal(keyIn)
	var keys []Keyword
	json.Unmarshal(keysJson, &keys)

	return keys
}

// Create new instance of keywordsCache
func NewKeyCache(cache Storage) *KeywordsCache {
	_, err := cache.GetV("_keys")
	isEmpty := false
	if err != nil {
		isEmpty = true
	}

	c := &KeywordsCache{
		cache: cache,
		key:   "_keys",
		next:  "_keys_next",
		isEmpty: isEmpty,
	}
	return c
}
