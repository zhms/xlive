package xutils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func HttpMakeGetUrl(url string, params map[string]interface{}) string {
	if params == nil {
		return url
	}
	if len(params) == 0 {
		return url
	}
	url += "?"
	for k, v := range params {
		url += fmt.Sprintf("%s=%v&", k, v)
	}
	url = url[:len(url)-1]
	return url
}

func HttpGet(url string) ([]byte, error) {
	header := map[string]string{}
	return HttpGetEx(url, header)
}

func HttpPost(url string, data interface{}) ([]byte, error) {
	header := map[string]string{}
	bytes, _ := json.Marshal(data)
	fmt.Println(string(bytes))
	return HttpPostEx(url, bytes, header)
}

func HttpGetEx(url string, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("http status code: %d", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func HttpPostEx(url string, data []byte, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Add(k, v)
	}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return nil, fmt.Errorf("http status code: %d", response.StatusCode)
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
