package provider

import (
	"bytes"
	"context"
	"crypto/tls"
	"fmt"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"io"
	"io/ioutil"
	"net/http"
)

// Client -
type Client struct {
	HostURL    string
	HTTPClient *http.Client
	Token      string
	Auth       AuthStruct
}

// AuthStruct -
type AuthStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthResponse -
type AuthResponse struct {
	UserID   int    `json:"user_id`
	Username string `json:"username`
	Token    string `json:"token"`
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func sendRequest(ctx context.Context, reqmethod string, c *Client, body []byte, apiUrl string, apiName string) string {
	tflog.Info(ctx, apiName+"===请求体==="+string(body)+"===")

	targetUrl := apiUrl

	req, _ := http.NewRequest(reqmethod, targetUrl, bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.SetBasicAuth(c.Auth.Username, c.Auth.Password)

	// 创建一个HTTP客户端并发送请求
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	respn, err := client.Do(req)
	if err != nil {
		tflog.Error(ctx, apiName+"===发送请求失败==="+err.Error())
		panic(apiName + "===发送请求失败===")
	}
	defer respn.Body.Close()

	body, err2 := io.ReadAll(respn.Body)
	if err2 != nil {
		tflog.Error(ctx, apiName+"===发送请求失败==="+err2.Error())
		panic(apiName + "===发送请求失败===")
	}

	tflog.Info(ctx, apiName+"===响应状态码==="+string(respn.Status)+"===")
	tflog.Info(ctx, apiName+"===响应体==="+string(body)+"===")

	//if strings.HasSuffix(respn.Status, "200") && strings.HasSuffix(respn.Status, "201") && strings.HasSuffix(respn.Status, "204") {
	//	tflog.Info(ctx, apiName+"===请求响应失败===")
	//	tflog.Info(ctx, apiName+"===响应状态码==="+string(respn.Status)+"===")
	//	tflog.Info(ctx, apiName+"===响应体==="+string(body)+"===")
	//} else {
	//	// 打印响应结果
	//	tflog.Info(ctx, apiName+"===响应状态码==="+string(respn.Status)+"===")
	//	tflog.Info(ctx, apiName+"===响应体==="+string(body)+"===")
	//}
	return string(body)
}
