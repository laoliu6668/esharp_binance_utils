package binance

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/laoliu6668/esharp_binance_utils/util"
)

func (c *ApiConfigModel) Get(gateway, path string, data map[string]any) (body []byte, resp *http.Response, err error) {
	return c.Request("GET", gateway, path, data, 0, true)
}

func (c *ApiConfigModel) Post(gateway, path string, data map[string]any) (body []byte, resp *http.Response, err error) {
	return c.Request("POST", gateway, path, data, 0, true)
}

func (c *ApiConfigModel) Delete(gateway, path string, data map[string]any) (body []byte, resp *http.Response, err error) {
	return c.Request("DELETE", gateway, path, data, 0, true)
}

// const proto = "http://"

const proto = "https://"

// 获取TRONSCAN API数据
func (c *ApiConfigModel) Request(method, gateway, path string, data map[string]any, timeout time.Duration, sign bool) (body []byte, resp *http.Response, err error) {

	if timeout == 0 {
		timeout = time.Second * 10
	}

	// 创建http client
	client := &http.Client{
		Timeout: timeout,
	}
	if UseProxy {
		uri, _ := url.Parse(fmt.Sprintf("http://%s", ProxyUrl))
		fmt.Printf("uri: %v\n", uri)
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(uri),
		}
	}
	if data == nil {
		data = make(map[string]any)
	}

	if sign {
		// 设置 时间戳
		data["timestamp"] = time.Now().UnixMilli()
		data["recvWindow"] = timeout.Milliseconds() // 同步安全时间
		// 签名
		data["signature"] = Signature(data, c.SecretKey)
	}

	// 构造query
	// url := proto + gateway + path

	// 声明 body
	// var reqBody io.Reader
	// if method == "POST" || method == "PUT" || method == "DELETE" {
	// 	if len(data) != 0 {
	// 		buf, _ := json.Marshal(data)
	// 		fmt.Printf("reqBody: %s\n", buf)
	// 		// 添加body
	// 		p := util.HttpBuildQuery(data)
	// 		reqBody = strings.NewReader(p)
	// 	}
	// 	// url = GetQueryUrl("https://", gateway, path, data)

	// } else if method == "GET" {
	// 	// url = GetQueryUrl(proto, gateway, path, data)
	// } else {
	// 	err = errors.New("不支持的http方法")
	// 	return
	// }
	url := GetQueryUrl(proto, gateway, path, data)
	// return nil, nil, errors.New("test")
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return
	}
	req.Header.Add("X-MBX-APIKEY", ApiConfig.AccessKey)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	resp, err = client.Do(req)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		msg := string(body)
		mp := map[string]any{}
		err = json.Unmarshal(body, &mp)
		if err == nil {
			if _, ok := mp["msg"]; ok {
				msg = fmt.Sprintf("%v", mp["msg"])
			}
		}
		return nil, nil, fmt.Errorf("http %v %v", resp.StatusCode, string(msg))
	}
	body, err = io.ReadAll(resp.Body)
	if err != nil {
		return
	}
	// fmt.Printf("body: %s\n", body)

	return
}

func GetQueryUrl(proto, gateway, path string, queryMap map[string]any) string {
	return fmt.Sprintf("%s%s%s?%s", proto, gateway, path, util.HttpBuildQuery(queryMap))

}

func Signature(args map[string]any, key string) string {
	str1 := util.HttpBuildQuery(args)
	h := hmac.New(sha256.New, []byte(key))
	h.Write([]byte(str1))
	return fmt.Sprintf("%x", h.Sum(nil))
	// return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

func UTCTimeNow() string {
	return time.Now().In(time.UTC).Format("2006-01-02T15:04:05")
}
