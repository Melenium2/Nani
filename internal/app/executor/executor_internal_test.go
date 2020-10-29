package executor

import (
	"Nani/internal/app/cache"
	"Nani/internal/app/config"
	"Nani/internal/app/inhuman"
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"sync"
	"testing"
	"time"
)

type mock_api struct {
}

func (m mock_api) App(bundle string) (*inhuman.App, error) {
	return &inhuman.App{
		Bundle: bundle,
	}, nil
}

func (m mock_api) Keys(title, description, shortDescription, reviews string) (inhuman.Keywords, error) {
	return map[string]int{
		"key3": 3,
		"key2": 2,
		"key":  1,
	}, nil
}

func (m mock_api) Flow(key string) ([]inhuman.App, error) {
	return []inhuman.App{
		{Bundle: "1"}, {Bundle: "2"}, {Bundle: "3"},
	}, nil
}

func (m mock_api) Endpoint(endpoint string) string {
	panic("implement me")
}

func (m mock_api) DevApps(devid string) ([]inhuman.App, error) {
	panic("implement me")
}

func (m mock_api) Request(endpoint, method string, data interface{}, response interface{}) error {
	panic("implement me")
}

type mock_storage struct {
	cache map[string]interface{}
	mutex sync.Mutex
}

func (m *mock_storage) Set(key string, value interface{}) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.cache[key] = value
}

func (m *mock_storage) GetV(key string) (interface{}, error) {
	//m.mutex.Lock()
	//defer m.mutex.Unlock()
	v, ok := m.cache[key]
	if !ok {
		return nil, fmt.Errorf("error key")
	}

	return v, nil
}

type mock_repo struct {
	Db map[int]*inhuman.App
}

func (m *mock_repo) Insert(ctx context.Context, app *inhuman.App) error {
	index := len(m.Db)
	m.Db[index+1] = app

	return nil
}

func (m *mock_repo) InsertBatch(ctx context.Context, apps []*inhuman.App) error {
	index := len(m.Db)
	for _, v := range apps {
		index += 1
		m.Db[index] = v
	}

	return nil
}

func TestDeclareTaskMock_ShouldLoadToCacheAllBundlesFromFile_NoError(t *testing.T) {
	ex := Executor{
		externalApi: mock_api{},
		cache:       &mock_storage{cache: make(map[string]interface{})},
	}

	err := ex.declareTask("../../../bundles.txt")
	assert.NoError(t, err)
	v, err := ex.cache.GetV("bundles")
	assert.NoError(t, err)
	assert.NotNil(t, v)
	keys := v.([]string)
	assert.Greater(t, len(keys), 0)

	b, err := ioutil.ReadFile("../../../bundles.txt")
	assert.NoError(t, err)
	l := strings.Split(string(b), "\n")
	assert.Equal(t, len(l), len(keys))
}

func TestStoreAppsMock_ShouldStoreAllAppsInADb_NoError(t *testing.T) {
	r := &mock_repo{Db: make(map[int]*inhuman.App)}
	ex := Executor{
		cache:       &mock_storage{cache: make(map[string]interface{})},
		externalApi: mock_api{},
		repository:  r,
		db:          make(databaseCh),
	}
	go ex.selector()
	apps := make([]string, 0)
	for i := 0; i < 51; i++ {
		apps = append(apps, "1")
	}
	ex.storeApps(false, apps...)
	time.Sleep(time.Second * 3)

	assert.Equal(t, 51, len(r.Db))
}

func TestStoreAppsMock_ShouldStartStoreKeywordFunctionAndStoreBothValues_NoError(t *testing.T) {
	c := &mock_storage{cache: make(map[string]interface{})}
	r := &mock_repo{Db: make(map[int]*inhuman.App)}
	ex := Executor{
		cache:       c,
		externalApi: mock_api{},
		keyCache:    cache.NewKeyCache(c),
		repository:  r,
		db:          make(databaseCh),
		config:      config.Config{KeysCount: 10},
	}
	go ex.selector()
	apps := make([]string, 0)
	for i := 0; i < 51; i++ {
		apps = append(apps, fmt.Sprintf("%d", i))
	}
	ex.storeApps(true, apps...)
	time.Sleep(time.Second * 3)

	for i := 0; i < 150; i++ {
		key, err := ex.keyCache.Next()
		assert.NoError(t, err)
		assert.NotEmpty(t, key)
	}

	v, err := ex.cache.GetV("last")
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "50", v.(string))
}

func TestStoreKeywordsMock_ShouldStoreKeywordsToCache_NoError(t *testing.T) {
	c := &mock_storage{cache: make(map[string]interface{})}
	ex := Executor{
		cache:       c,
		externalApi: mock_api{},
		keyCache:    cache.NewKeyCache(c),
		db:          make(databaseCh),
		config:      config.Config{KeysCount: 10},
	}
	go ex.selector()
	ex.storeKeywords(&inhuman.App{})
	time.Sleep(time.Second * 3)

	for i := 0; i < 3; i++ {
		key, err := ex.keyCache.Next()
		assert.NoError(t, err)
		assert.NotEmpty(t, key)
	}
}

func TestSaveErrorMock_ShouldSaveNewErrorToCache_NoError(t *testing.T) {
	ex := Executor{
		cache: &mock_storage{cache: make(map[string]interface{})},
	}

	ex.saveError("bundle", "com.1", fmt.Errorf("error1"))
	ex.saveError("bundle", "com.2", fmt.Errorf("error2"))
	ex.saveError("keys", "keyword", fmt.Errorf("error1"))

	e, err := ex.cache.GetV("errors")
	assert.NoError(t, err)
	assert.NotNil(t, e)
	ers := e.([]ExecutorError)

	assert.Equal(t, "bundle", ers[0].t)
	assert.Equal(t, "bundle", ers[1].t)
	assert.Equal(t, "keys", ers[2].t)

	assert.Equal(t, "com.1", ers[0].bundle)
	assert.Equal(t, "com.2", ers[1].bundle)
	assert.Equal(t, "keyword", ers[2].bundle)
}

func TestSelectorMock_ShouldSaveAppsIfApplicationCanceled_NoError(t *testing.T) {
	r := &mock_repo{Db: make(map[int]*inhuman.App)}
	ex := Executor{
		db:         make(databaseCh),
		config:     config.Config{KeysCount: 10},
		repository: r,
	}
	go ex.selector()

	for i := 0; i < 20; i++ {
		ex.db <- &inhuman.App{Bundle: "1"}
	}

	ex.cancel = true

	time.Sleep(time.Second * 3)

	assert.Equal(t, 20, len(r.Db))
}

func TestAppsBatchMock_ShouldStoreApplicationDataFromMainPage_NoError(t *testing.T) {
	ctx := context.Background()
	r := &mock_repo{Db: make(map[int]*inhuman.App)}
	c := &mock_storage{cache: make(map[string]interface{})}
	ex := Executor{
		cache:       c,
		db:          make(databaseCh),
		config:      config.Config{KeysCount: 10},
		repository:  r,
		externalApi: mock_api{},
		keyCache:    cache.NewKeyCache(c),
		ctx:         ctx,
		wait:        make(chan struct{}),
	}
	go ex.selector()

	ex.keyCache.Set("123")
	ex.keyCache.Set("123")
	ex.keyCache.Set("123")
	go func() {
		time.Sleep(time.Second * 3)
		ex.cancel = true
	}()
	ex.appsBatch()

	time.Sleep(time.Second * 3)
	assert.Equal(t, 9, len(r.Db))

	k, err := ex.keyCache.Next()
	assert.Error(t, err)
	assert.Empty(t, k)
}
