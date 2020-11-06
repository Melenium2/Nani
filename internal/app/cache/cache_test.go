package cache_test

import (
	"Nani/internal/app/cache"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestSetV_ShouldStoreToCacheNewValue_NoErrors(t *testing.T) {
	k, v := "key", "value"
	c := cache.New(true)
	c.Set(k, v)
	res, err := c.GetV(k)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, v, res.(string))
}

func TestGetV_ShouldReturnValueFromCache_NoErrors(t *testing.T) {
	k, v := "key", "value"
	c := cache.New(true)
	c.Set(k, v)
	res, err := c.GetV(k)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, v, res.(string))
}

func TestGetV_ShouldReturnErrorCozKeyIsInvalid_Error(t *testing.T) {
	k, v := "key", "value"
	c := cache.New(true)
	c.Set(k, v)
	res, err := c.GetV(k + "1")
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestLoad_ShouldStoreAllJsonToCache_NoError(t *testing.T) {
	_, err := ioutil.ReadFile("./cache.json")
	if err != nil {
		t.Skipf("file /cache.json not found %s", err)
	}

	c := cache.New(false)
	v, err := c.GetV("key")
	assert.NoError(t, err)
	assert.NotNil(t, v)
	assert.Equal(t, "value", v.(string))
}

func TestLoad_ShouldThrowPanicCozInalidJosn_Error(t *testing.T) {
	assert.Panics(t, func() {
		j, _ := json.Marshal(`{"123" "1321323}`)
		ioutil.WriteFile("./cache.json", j, 0644)
		cache.New(false)
	})
	os.Remove("./cache.json")
}

func TestLoad_ShouldReturnCorrectErrorsObject_Error(t *testing.T) {
	f, _ := ioutil.ReadFile("./cache.json")
	var m map[string]cache.Item
	json.Unmarshal(f, &m)
	e, _ := json.Marshal(m["_keys"].V)
	var l []cache.Keyword
	r := json.Unmarshal(e, &l)
	t.Log(e)
	t.Log(r)
	t.Log(l)
	//c := cache.New(false)
	//v, err := c.GetV("errors")
	//assert.NoError(t, err)
	//assert.NotNil(t, v)
	//t.Logf("%v", v)
	//e, ok := v.([]executor.ExecutorError)
	//assert.True(t, ok)
	//assert.NotNil(t, e)
}

func TestDump_ShouldSaveValidDump_NoError(t *testing.T) {
	k, v := "key", "value"
	c := cache.New(true)
	c.Set(k, v)
	res, err := c.GetV(k)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, v, res.(string))
	c.Dump()

	f, err := ioutil.ReadFile("./cache.json")
	assert.NoError(t, err)

	var m map[string]cache.Item
	json.Unmarshal(f, &m)

	assert.NotNil(t, m)
	assert.Equal(t, "value", m["key"].V)

	os.Remove("./cache.json")
}




