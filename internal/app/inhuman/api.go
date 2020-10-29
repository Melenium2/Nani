package inhuman

import (
	"Nani/internal/app/config"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)


// Enterprise Api Service decorator
type ExternalApi interface {
	App(bundle string) (*App, error)
	Keys(title, description, shortDescription, reviews string) (Keywords, error)
	Flow(key string) ([]App, error)
	Endpoint(endpoint string) string
	Request(endpoint, method string, data interface{}, response interface{}) error
}

// Api for interaction with external service (Enterprise Api Service)
type InhumanApi struct {
	config *config.Config
}

// Get application info by bundle from external api
// @bundle: string (application bundle for example 'com.123.pepega')
// @return @App: *App (application struct), @Error: error (error if request was not successful)
func (api *InhumanApi) App(bundle string) (*App, error) {
	var app *App
	err := api.Request(api.Endpoint("bundle"), "post", map[string]string{
		"query": bundle,
		"hl": api.config.Hl,
		"gl": api.config.Gl,
	}, &app)

	if err != nil {
		return nil, err
	}

	return app, nil
}

// Method call EAS api and return keywords with their weight from application text
// @Title, @Description, @shortDescription, @reviews: String (text fields from application)
// @return Keywords: map[string]int (with keywords and their weight), @error: error
func (api *InhumanApi) Keys(title, description, shortDescription, reviews string) (Keywords, error) {
	keywords := make(Keywords)
	err := api.Request(api.Endpoint("keywords_from"), "post", map[string]interface{}{
		"title": title,
		"description": description,
		"shortDescription": shortDescription,
		"reviews": reviews,
		"keysCount": api.config.KeysCount,
		"lang": api.config.Hl,
	}, &keywords)

	if err != nil {
		return nil, err
	}

	return keywords, nil
}

// Method call EAS api and return top N application from google play main page
// @key: string (keywords for search)
// @return []App (list of application metainfo) @error Error
func (api *InhumanApi) Flow(key string) ([]App, error) {
	apps := make([]App, 0)
	err := api.Request(api.Endpoint("mainPage"), "post", map[string]interface{} {
		"query": key,
		"hl": api.config.Hl,
		"gl": api.config.Gl,
		"count": api.config.AppsCount,
	}, &apps)

	if err != nil {
		return nil, err
	}

	return apps, nil
}

// Generate url to EAS api from string
// @endpoint: string (endpoint name)
// @return: string (full endpoint url)
func (api *InhumanApi) Endpoint(endpoint string) string {
	return fmt.Sprintf("%s/%s", api.config.ApiUrl, endpoint)
}

//Make request to EAS to the given endpoint
//@endpoint: String (api endpoint)
//@method: String (request method)
//@data: interface (any post object)
//@response: interface (any response object pointer)
//@return error
func (api *InhumanApi) Request(endpoint, method string, data interface{}, response interface{}) error {
	var err error
	var b []byte
	b, err = json.Marshal(data)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(strings.ToUpper(method), endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", api.config.Key)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode > 200 {
		return fmt.Errorf("external api response with status %d", resp.StatusCode)
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return err
	}

	return nil
}

func New(config *config.Config) *InhumanApi {
	return &InhumanApi{
		config: config,
	}
}
