package apis

import (
	"encoding/json"
	"fmt"

	root "github.com/laoliu6668/esharp_binance_utils"
)

type ListenKey struct {
	ListenKey string `json:"listenKey"`
}

// 获取现货账户的listenKey
func GetSpotAccountListenKey() (data string, err error) {
	const flag = "Binance GetUserDataListenKey"
	body, _, err := root.ApiConfig.Request("POST", gateway_binance, "/api/v3/userDataStream", nil, 0, false)
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}
	res := ListenKey{}
	json.Unmarshal(body, &res)
	return res.ListenKey, nil
}

// 延长现货账户的listenKey
func KeepSpotAccountListenKey(listenKey string) (err error) {
	const flag = "Binance KeepSpotAccountListenKey"
	_, _, err = root.ApiConfig.Request("PUT", gateway_binance, "/api/v3/userDataStream", map[string]any{
		"listenKey": listenKey,
	}, 0, false)
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}
	return nil
}

// 获取期货账户的listenKey
func GetSwapAccountListenKey() (data string, err error) {
	const flag = "Binance GetSwapAccountListenKey"
	body, _, err := root.ApiConfig.Request("POST", gateway_fapi, "/fapi/v1/listenKey", nil, 0, false)
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}
	fmt.Printf("body: %s\n", body)
	res := ListenKey{}
	json.Unmarshal(body, &res)
	fmt.Printf("res: %v\n", res)
	return res.ListenKey, nil
}

// 获取期货账户的listenKey
func KeepSwapAccountListenKey(listenKey string) (err error) {
	const flag = "Binance KeepSwapAccountListenKey"
	_, _, err = root.ApiConfig.Request("PUT", gateway_fapi, "/fapi/v1/listenKey", map[string]any{
		"listenKey": listenKey,
	}, 0, false)
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}
	return nil
}
