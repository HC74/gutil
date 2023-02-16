package gutil

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

type M = map[string]string

type HttpClient struct {
	client *http.Client
}

// NewHttpClient 创建httpClient请求实例
func NewHttpClient() *HttpClient {
	return &HttpClient{
		client: &http.Client{},
	}
}

// httpRequest 通用请求方法
func httpRequest(c *http.Client, requestType, url string, headers M, dataBody io.Reader) ([]byte, error) {
	req, err := http.NewRequest(requestType, url, dataBody)
	if err != nil {
		return nil, err
	}

	// 设置请求头
	for key, value := range headers {
		req.Header.Set(key, value)
	}
	// post默认采用json
	if strings.EqualFold(requestType, "POST") {
		req.Header.Set("content-type", "application/json;charset=UTF-8")
	}

	// 发送请求
	resp, err := c.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, err
}

// httpRequest 发送get请求
func (c *HttpClient) httpRequest(requestType, url string, headers M, dataBody interface{}) ([]byte, error) {
	var dataJson = ""
	if strings.EqualFold(requestType, "POST") || dataBody != nil {
		jsonData, _ := json.Marshal(dataBody)
		dataJson = string(jsonData)
	}
	return httpRequest(c.client, requestType, url, headers, strings.NewReader(dataJson))
}

// GetStringN 发送get请求且不包含请求头
func (c *HttpClient) GetStringN(url string) (string, error) {
	return c.GetString(url, map[string]string{})
}

// PostStringN 发送post请求且不包含请求头
func (c *HttpClient) PostStringN(url string, body interface{}) (string, error) {
	return c.PostString(url, map[string]string{}, body)
}

// PostBodyN 发送post请求并不包含请求头且完成数据获取
func (c *HttpClient) PostBodyN(url string, body interface{}, data interface{}) error {
	return c.PostBody(url, map[string]string{}, body, data)
}

// GetBodyN 发送get请求并完成数据获取
func (c *HttpClient) GetBodyN(url string, data interface{}) error {
	return c.GetBody(url, map[string]string{}, data)
}

// GetString 发送get请求且包含请求头
func (c *HttpClient) GetString(url string, headers M) (string, error) {
	body, err := c.httpRequest("GET", url, headers, nil)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// GetBody 发送get请求并完成数据获取
func (c *HttpClient) GetBody(url string, headers M, data interface{}) error {
	jsonData, err := c.httpRequest("GET", url, headers, nil)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return err
	}
	return nil
}

// PostString 发送post请求返回string
func (c *HttpClient) PostString(url string, headers M, body interface{}) (string, error) {
	respBody, err := c.httpRequest("POST", url, headers, body)
	if err != nil {
		return "", err
	}
	return string(respBody), nil
}

// PostBody 发送post请求并完成数据获取
func (c *HttpClient) PostBody(url string, headers M, body interface{}, data interface{}) error {
	jsonData, err := c.httpRequest("POST", url, headers, body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return err
	}
	return nil
}

// == 泛型

// GetBody 泛型Get
func GetBody[T any](url string, headers M) (*T, error) {
	c := &http.Client{}
	var data *T
	req, err := httpRequest(c, "GET", url, headers, nil)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(req, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}

// GetPost 泛型post
func GetPost[T any](url string, headers M, body interface{}) (*T, error) {
	c := &http.Client{}
	var dataJson string
	if body != nil {
		jsonData, _ := json.Marshal(body)
		dataJson = string(jsonData)
	}
	req, err := httpRequest(c, "POST", url, headers, strings.NewReader(dataJson))
	if err != nil {
		return nil, err
	}
	var data *T
	err = json.Unmarshal(req, &data)
	if err != nil {
		return nil, err
	}
	return data, nil
}
