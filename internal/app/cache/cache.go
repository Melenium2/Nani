package cache

import (
	"Nani/internal/app/file"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path"
	"sync"
	"syscall"
)

var debug bool = false

type Storage interface {
	Set(key string, value interface{})
	GetV(key string) (interface{}, error)
}

type Cache struct {
	Clear     bool
	Debug     bool
	cachename string
	store     map[string]Item
	mutex     sync.Mutex
}

func (c *Cache) Set(key string, value interface{}) {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.store[key] = Item{value, 0}
}

func (c *Cache) GetV(key string) (interface{}, error) {
	v, ok := c.store[key]
	if !ok {
		return nil, errors.New(fmt.Sprintf("value with key %s not found", key))
	}

	return v.V, nil
}

func (c *Cache) dump() {
	if debug {
		log.Print("start dump")
	}
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("creating cache dump coz %s", <-sig)
		if debug {
			log.Print("sig complete")
		}

		close(sig)
		if len(c.store) > 0 {
			f := file.New(c.cachename)
			str, err := json.Marshal(c.store)
			if err != nil {
				panic(err)
			}
			f.WriteLines(string(str))

			log.Print("cache crated")
		} else {
			log.Print("skip creating")

			return
		}
	}()
}

func (c *Cache) load() {
	if c.Clear {
		err := os.Remove(c.cachename)
		if err != nil {
			log.Printf("Can not remove file %s. File not found", c.cachename)
		}
	} else {
		f := file.New(c.cachename)
		b, err := f.ReadAll()
		if err != nil || len(b) == 0 {
			log.Print("File not exist")
			return
		}

		err = json.Unmarshal(b, &c.store)
		if err != nil {
			panic(err)
		}
	}
}

func New(new bool, cachefile ...string) *Cache {
	filename := path.Join("./cache.json")
	if len(cachefile) > 0 {
		filename = cachefile[0]
	}

	c := &Cache{
		cachename: filename,
		store:     make(map[string]Item),
		Clear:     new,
	}
	c.load()
	c.dump()
	return c
}
