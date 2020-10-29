package inhuman_test

import (
	"Nani/internal/app/config"
	"Nani/internal/app/inhuman"
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

type inhuman_api_mock struct {
	ExpectedCode int
	ExpectedBody interface{}
}

func (m *inhuman_api_mock) Flow(key string) ([]inhuman.App, error) {
	apps := make([]inhuman.App, 0)
	err := m.Request(m.Endpoint("mainPage"), "post", map[string]interface{} {
		"query": key,
		"hl": "13",
		"gl": "123",
		"count": 50,
	}, &apps)

	if err != nil {
		return nil, err
	}

	return apps, nil
}

func (m *inhuman_api_mock) Keys(title, description, shortDescription, reviews string) (inhuman.Keywords, error) {
	keywords := make(inhuman.Keywords)
	err := m.Request(m.Endpoint("keywords_from"), "post", map[string]string{
		"title":            title,
		"description":      description,
		"shortDescription": shortDescription,
		"reviews":          reviews,
	}, &keywords)

	if err != nil {
		return nil, err
	}

	return keywords, nil
}

func (m *inhuman_api_mock) App(bundle string) (*inhuman.App, error) {
	var app *inhuman.App
	err := m.Request(m.Endpoint("bundle"), "post", map[string]string{
		"query": bundle,
		"hl":    "en",
		"gl":    "us",
	}, &app)

	if err != nil {
		return nil, err
	}

	return app, nil
}

func (m *inhuman_api_mock) Endpoint(endpoint string) string {
	return fmt.Sprintf("/%s", endpoint)
}

func (m *inhuman_api_mock) Request(endpoint, method string, data interface{}, response interface{}) error {
	var err error
	var b []byte
	b, err = json.Marshal(data)
	if err != nil {
		return err
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		app := m.ExpectedBody
		b, _ := json.Marshal(app)
		w.WriteHeader(m.ExpectedCode)
		io.WriteString(w, string(b))
	}

	r := httptest.NewRequest(method, endpoint, bytes.NewReader(b))
	w := httptest.NewRecorder()
	handler(w, r)

	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		return err
	}

	return nil
}

func (m *inhuman_api_mock) DevApps(devid string) ([]inhuman.App, error) {
	apps := make([]inhuman.App, 0)
	err := m.Request(m.Endpoint("devapps"), "post", map[string]interface{} {
		"query": devid,
		"hl":    "en",
		"gl":    "us",
		"count": 100,
	}, &apps)

	if err != nil {
		return nil, err
	}

	return apps, nil
}

/*
Successful test
============================================================================================
 */

func TestApp_ShouldMakeRequestToExternalApiAndWriteToAppResult_NoErrors(t *testing.T) {
	api := &inhuman_api_mock{
		ExpectedCode: 200,
		ExpectedBody: &inhuman.App{
			Bundle:     "123",
			Categories: "GAME",
		},
	}
	res, err := api.App("exmaple")
	assert.NoError(t, err)
	assert.Equal(t, "123", res.Bundle)
	assert.Equal(t, "GAME", res.Categories)
}

func TestEndpoint_ShouldReturnCorrectEndpointString_NoErrors(t *testing.T) {
	api := &inhuman_api_mock{}
	res := api.Endpoint("bundle")
	assert.Equal(t, "/bundle", res)
}

func TestRequest_ShouldMakeRequestToExternalApi_NoErrors(t *testing.T) {
	api := &inhuman_api_mock{}
	err := api.Request("/exmaple", "get", struct{}{}, struct{}{})
	assert.NoError(t, err)
}

func TestKeys_ShouldReturnObjectWithKeywords_NoErrors(t *testing.T) {
	api := &inhuman_api_mock{
		ExpectedCode: 200,
		ExpectedBody: inhuman.Keywords{
			"key": 1,
			"key2": 2,
			"key3": 3,
		},
	}
	keywords, err := api.Keys("1", "2", "3", "4")
	assert.NoError(t, err)
	assert.Equal(t, 3, len(keywords))
}

func TestFlow_ShouldReturnListWithApplications_NoError(t *testing.T) {
	api := &inhuman_api_mock{
		ExpectedCode: 200,
		ExpectedBody: []inhuman.App {
			inhuman.App{ Bundle: "123"}, { Bundle: "222"},
		},
	}
	apps, err := api.Flow("car")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(apps))
}

func TestDevApps_ShouldRetrunListOfDeveloperApps_NoError(t *testing.T) {
	api := &inhuman_api_mock{
		ExpectedCode: 200,
		ExpectedBody: []inhuman.App {
			{ Bundle: "123"}, {Bundle: "qwe"},
		},
	}
	apps, err := api.DevApps("228")
	assert.NoError(t, err)
	assert.Equal(t, 2, len(apps))
}

/*
===================================================================================================
 */

type inhuman_api_mock_fail struct {
	ExpectedCode         int
	ExpectedResponseBody interface{}
}

func (m *inhuman_api_mock_fail) Flow(key string) ([]inhuman.App, error) {
	apps := make([]inhuman.App, 0)
	err := m.Request(m.Endpoint("mainPage"), "post", map[string]interface{} {
		"query": key,
		"hl": "13",
		"gl": "123",
		"count": 50,
	}, &apps)

	if err != nil {
		return nil, err
	}

	return apps, nil
}

func (m *inhuman_api_mock_fail) Keys(title, description, shortDescription, reviews string) (inhuman.Keywords, error) {
	keywords := make(inhuman.Keywords)
	err := m.Request(m.Endpoint("keywords_from"), "post", map[string]string{
		"title":            title,
		"description":      description,
		"shortDescription": shortDescription,
		"reviews":          reviews,
	}, &keywords)

	if err != nil {
		return nil, err
	}

	return keywords, nil
}

func (m *inhuman_api_mock_fail) App(bundle string) (*inhuman.App, error) {
	return nil, m.Request(m.Endpoint("bundle"), "post", struct{}{}, struct{}{})
}

func (m *inhuman_api_mock_fail) Endpoint(endpoint string) string {
	return fmt.Sprintf("/%s", endpoint)
}

func (m *inhuman_api_mock_fail) Request(endpoint, method string, data interface{}, response interface{}) error {
	var err error
	var b []byte
	b, err = json.Marshal(data)
	if err != nil {
		return err
	}

	handler := func(w http.ResponseWriter, r *http.Request) {
		b, _ := json.Marshal(m.ExpectedResponseBody)
		w.WriteHeader(m.ExpectedCode)
		io.WriteString(w, string(b))
	}

	r := httptest.NewRequest(method, endpoint, bytes.NewReader(b))
	w := httptest.NewRecorder()
	handler(w, r)

	if w.Code >= 300 {
		return errors.New("response with fail status")
	}

	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		return err
	}

	return nil
}

func (m *inhuman_api_mock_fail) DevApps(devid string) ([]inhuman.App, error) {
	apps := make([]inhuman.App, 0)
	err := m.Request(m.Endpoint("devapps"), "post", map[string]interface{} {
		"query": devid,
		"hl":    "en",
		"gl":    "us",
		"count": 100,
	}, &apps)

	if err != nil {
		return nil, err
	}

	return apps, nil
}

/*
Unsuccessful test
============================================================================================
*/

func TestApp_ShouldReturn500ErrorOrJsonMarshalError_Error(t *testing.T) {
	api := &inhuman_api_mock_fail{
		ExpectedCode: 500,
	}
	res, err := api.App("bundle")
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestEndpoint_ShouldReturn500Error_Error(t *testing.T) {
	api := &inhuman_api_mock_fail{
		ExpectedCode: 500,
		ExpectedResponseBody: map[string]string{
			"bundle": "app",
		},
	}
	err := api.Request("/bundle", "get", struct{}{}, struct{}{})
	assert.Error(t, err)
	assert.Equal(t, "response with fail status", err.Error())
}

func TestEndpoint_ShouldReturnUnpredictableBody_Error(t *testing.T) {
	api := &inhuman_api_mock_fail{
		ExpectedCode:         200,
		ExpectedResponseBody: `{ "bundle": "bundle" }]`,
	}
	err := api.Request("/bundle", "get", struct{}{}, &inhuman.App{})
	assert.Error(t, err)
	assert.Equal(t, "json: cannot unmarshal string into Go value of type inhuman.App", err.Error())
}

func TestEndpoint_ShouldReturnErrorCozIncorrectData_Error(t *testing.T) {
	api := &inhuman_api_mock_fail{
		ExpectedCode:         200,
		ExpectedResponseBody: `{ "bundle": "bundle" }`,
	}
	err := api.Request("/bundle", "get", make(chan int), &inhuman.App{})
	assert.Error(t, err)
}

func TestKeys_ShouldReturnWrongResultFromRequest_Error(t *testing.T) {
	api := &inhuman_api_mock_fail{
		ExpectedCode:         200,
		ExpectedResponseBody: `[{ "1": "1" }]`,
	}
	keywords, err := api.Keys("/bundle", "get", "1", "3")
	assert.Error(t, err)
	assert.Equal(t, 0, len(keywords))
}

func TestFlow_ShouldReturnWrongResultFromRequestIncorrectDataStruct_Error(t *testing.T) {
	api := &inhuman_api_mock_fail{
		ExpectedCode:         200,
		ExpectedResponseBody: `[{ "1": "1" }]`,
	}
	apps, err := api.Flow("car")
	assert.Error(t, err)
	assert.Nil(t, apps)
}

func TestFlow_ShouldReturnWrongResultFromRequestIncorrectDataArray_Error(t *testing.T) {
	api := &inhuman_api_mock_fail{
		ExpectedCode:         200,
		ExpectedResponseBody: `{ "1": "1" }`,
	}
	apps, err := api.Flow("car")
	assert.Error(t, err)
	assert.Nil(t, apps)
}

func TestDevApps_ShouldReturnErrorCode_Error(t *testing.T) {
	api := &inhuman_api_mock_fail{
		ExpectedCode:         500,
		ExpectedResponseBody: `{ "1": "1" }`,
	}
	apps, err := api.DevApps("car")
	assert.Error(t, err)
	assert.Nil(t, apps)
}

func TestDevApps_ShouldReturnWrongResult_Error(t *testing.T) {
	api := &inhuman_api_mock_fail{
		ExpectedCode:         200,
		ExpectedResponseBody: `{ "1": "1" }""`,
	}
	apps, err := api.DevApps("car")
	assert.Error(t, err)
	assert.Nil(t, apps)
}

func TestDevApps_ShouldReturnErrorCozIncorrectData_Error(t *testing.T) {
	api := &inhuman_api_mock_fail{
		ExpectedCode:         404,
		ExpectedResponseBody: `{ "1": "1" }`,
	}
	apps, err := api.DevApps("")
	assert.Error(t, err)
	assert.Nil(t, apps)
}

/*
=============================================================================
 */

func Config() *config.Config {
	os.Setenv("api_key", "Security 3923cf9a417e73be95b40dc5db60c97dcb876a61")
	c := config.New()
	c.AppsCount = 250
	c.KeysCount = 10

	return c
}

var bundle = "com.and.wareternal"

func TestEndpoint_ShouldConcatDefaultUrlAndEndpoint_NoError(t *testing.T) {
	c := Config()
	api := inhuman.New(c)
	res := api.Endpoint("bundle")
	assert.Equal(t, c.ApiUrl + "/bundle", res)
}

func TestApp_ShouldReturnAppInformationFromApi_NoError(t *testing.T) {
	api := inhuman.New(Config())
	res, err := api.App(bundle)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, bundle, res.Bundle)
}

func TestKeys_ShouldReturnKeywordsFromGivenText_NoError(t *testing.T) {
	api := inhuman.New(Config())
	res, err := api.App(bundle)
	assert.NoError(t, err)
	assert.NotNil(t, res)

	assert.NotEmpty(t, res.Title)
	assert.NotEmpty(t, res.Description)
	assert.NotEmpty(t, res.ShortDescription)

	keys, err := api.Keys(res.Title, res.Description, res.ShortDescription, "")
	assert.Greater(t, len(keys), 0)
}

func TestFlow_ShouldReturnMainPageApps_NoError(t *testing.T) {
	c := Config()
	api := inhuman.New(c)
	res, err := api.Flow("car")
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Greater(t, len(res), 0)
	assert.Greater(t, len(res), c.AppsCount - 30)
}

func TestDevApps_ShouldReturnAllAppsApplications_NoError(t *testing.T) {
	api := inhuman.New(Config())
	res, err := api.App(bundle)
	assert.NoError(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, bundle, res.Bundle)

	apps, err := api.DevApps(res.DeveloperId)
	assert.NoError(t, err)
	assert.Greater(t, len(apps), 0)
}

func TestFlow_ShouldReturnAppsFor10Keys_NoError(t *testing.T) {
	ti := time.Now()
	c := Config()
	api := inhuman.New(c)
	keys := []string {"car", "cart", "car games", "game for kids", "russian mobiles", "anime", "anime games", "wallpapers", "key", "door"}
	for _, k := range keys {
		res, err := api.Flow(k)
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Greater(t, len(res), 0)
		t.Log(k)
	}
	t.Log(time.Now().Sub(ti).Seconds() * 6000)
}

func TestApp_ShouldReturnErrorIfBundleIsWrong_Error(t *testing.T) {
	api := inhuman.New(Config())
	res, err := api.App("")
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestApp_ShouldReturnErrorCozKeyIsIncorrect_Error(t *testing.T) {
	api := inhuman.New(Config())
	res, err := api.App("dfghadsvadkasdasdskjdsnkjdna123ad;lmsakda")
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestDevApps_ShouldReturnErrorCozDevIdIsIncorrect_Error(t *testing.T) {
	api := inhuman.New(Config())
	res, err := api.DevApps("dfghadsvadkasdasdskjdsnkjdna123ad;lmsakda")
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestDevApps_ShouldReturnErrorCozDevIdIsEmpty_Error(t *testing.T) {
	api := inhuman.New(Config())
	res, err := api.DevApps("")
	assert.Error(t, err)
	assert.Nil(t, res)
}
