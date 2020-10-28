package cache

import "errors"

type KeyStorage interface {
	Set(key string) error
	Next() (string, error)
	Rollback() error
}

type KeywordsCache struct {
	cache Storage
	key   string
	next  string
}

func (kc *KeywordsCache) Set(key string) error {
	keysIn, err := kc.cache.GetV(kc.key)
	keys, ok := keysIn.([]Keyword)
	if err != nil || !ok {
		kc.cache.Set(kc.key, []Keyword{{Pos: len(keys) + 1, Key: key}})
		return nil
	}

	kc.cache.Set(kc.key, append(keys, Keyword{Pos: len(keys) + 1, Key: key}))
	return nil
}

func (kc *KeywordsCache) Next() (string, error) {
	keyIn, err := kc.cache.GetV(kc.next)
	if err != nil {
		kc.cache.Set(kc.next, 0)
		k, err := kc.get(0)
		if err != nil {
			return "", err
		}
		return k.Key, nil
	}
	key := keyIn.(int) + 1


	keyword, err := kc.get(key)
	if err != nil {
		return "", errors.New("keywords are out of range")
	}

	kc.cache.Set(kc.next, key)

	return keyword.Key, nil
}

func (kc *KeywordsCache) Rollback() error {
	keyIn, err := kc.cache.GetV(kc.next)
	if err != nil {
		return nil
	}
	key := keyIn.(int)

	if key == 0 {
		return nil
	}

	kc.cache.Set(kc.next, key - 1)
	return nil
}

func (kc *KeywordsCache) get(i int) (Keyword, error) {
	keyIn, err := kc.cache.GetV(kc.key)
	if err != nil {
		return Keyword{}, errors.New("keywords cache are empty")
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

func NewKeyCache(cache Storage) *KeywordsCache {
	c :=  &KeywordsCache{
		cache: cache,
		key:   "_keys",
		next:  "_keys_next",
	}
	return c
}
