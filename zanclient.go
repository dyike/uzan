package uzan

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	// OpenAPI open api url.
	OpenAPI = "https://open.youzan.com/api"
	// TimeFormatter time formatter.
	TimeFormatter = "2006-01-02 15:04:05"
)

// ZanClient zan client for api client.
type ZanClient struct {
	ClientID     string
	ClientSecret string
	KdtID        int64
	AccessToken  string
}

type stringer interface {
	String() string
}

func getMd5String(s string) string {
	h := md5.New()
	h.Write([]byte(s))
	return hex.EncodeToString(h.Sum(nil))
}

func getKeyValue(key interface{}) string {
	var ret = ""
	switch valueType := key.(type) {
	case string:
		ret = valueType
	case stringer:
		ret = valueType.String()
	default:
		ret = "Error in value connert"
	}
	return ret
}

// getURL get url by youzan method and version.
func (c *ZanClient) getURL(method string, version string) string {
	methodArray := strings.Split(method, ".")
	// 请求action
	action := methodArray[len(methodArray)-1]
	service := strings.Join(methodArray[0:len(methodArray)-1], ".")
	httpURL := "/" + service + "/" + version + "/" + action
	return httpURL
}

// ZanRequset zan client request.
func (c *ZanClient) ZanRequset(
	method string,
	params map[string]interface{},
	version string,
	files map[string]interface{}) ([]byte, error) {
	url := OpenAPI + "/oauthentry"
	service := c.getURL(method, version)
	var paramsMap map[string]interface{}

	paramsMap = params
	paramsMap["access_token"] = c.AccessToken

	url += service

	resp, err := c.sendRequest(url, "post", paramsMap, files)
	defer resp.Body.Close()

	if err == nil {
		if resp.StatusCode != http.StatusOK {
			err = errors.New("http error code: " + string(resp.StatusCode) + " reason: " + resp.Status)
		}
	}

	var result []byte
	if err != nil {
		result, err = ioutil.ReadAll(resp.Body)
	}

	return result, err

}

func (c *ZanClient) sendRequest(rawURL string, method string, params map[string]interface{}, files map[string]interface{}) (*http.Response, error) {
	httpClient := &http.Client{}

	var req *http.Request
	var err error
	if "GET" == strings.ToUpper(method) {
		httpURL := rawURL
		if len(params) > 0 {
			httpURL += "?"
			for key, value := range params {
				httpURL += url.QueryEscape(key)
				httpURL += "="
				httpURL += url.QueryEscape(getKeyValue(value))
				httpURL += "&"
			}
			httpURL = strings.TrimRight(httpURL, "&")
		}
		req, err = http.NewRequest("GET", httpURL, nil)
	} else if "POST" == strings.ToUpper(method) {
		jsonData, _ := json.Marshal(params)
		dataReader := bytes.NewReader(jsonData)
		req, err = http.NewRequest("POST", rawURL, dataReader)
		req.Header.Set("Content-Type", "application/json")
	} else {
		panic(errors.New("not support method"))
	}

	if err != nil {
		panic(err)
	}

	req.Header.Add("User-Agent", "X-UZan-Client")
	return httpClient.Do(req)
}
