package apis

import (
	"encoding/json"
	"fmt"
	"strings"

	root "github.com/laoliu6668/esharp_binance_utils"
	"github.com/laoliu6668/esharp_binance_utils/util"
)

const gateway_vision = "data-api.binance.vision"

type ApiResponseFiliter struct {
	Data []map[string]any `json:"symbols"`
}

// 获取全局过滤器
// https://developers.binance.com/docs/zh-CN/binance-spot-api-docs/filters
func GetFiliters() (data []map[string]any, err error) {
	const flag = "binance GetFiliters"
	body, _, err := root.ApiConfig.Request("GET", gateway_vision, "/api/v3/exchangeInfo", nil, 0, false)
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}
	// util.WriteTestJsonFile(flag, body)
	res := ApiResponseFiliter{}
	d := json.NewDecoder(strings.NewReader(string(body)))
	d.UseNumber()
	err = d.Decode(&res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}
	// 过滤
	data = []map[string]any{}
	for _, v := range res.Data {
		if v["quoteAsset"] == "USDT" && v["status"] == "TRADING" {
			data = append(data, v)
		}
	}
	return data, nil
}

// 获取期货交易对配置
// https://developers.binance.com/docs/zh-CN/binance-spot-api-docs/filters
func GetSwapSymbols() (data []map[string]any, err error) {
	const flag = "binance GetSwapSymbols"
	body, _, err := root.ApiConfig.Request("GET", gateway_fapi, "/fapi/v1/exchangeInfo", nil, 0, false)
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}
	util.WriteTestJsonFile(flag, body)
	res := ApiResponseFiliter{}
	d := json.NewDecoder(strings.NewReader(string(body)))
	d.UseNumber()
	err = d.Decode(&res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}
	// 过滤
	data = []map[string]any{}
	for _, v := range res.Data {
		if v["quoteAsset"] == "USDT" && v["status"] == "TRADING" {
			data = append(data, v)
		}
	}
	return data, nil
}
