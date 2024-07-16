package apis

import (
	"encoding/json"
	"fmt"

	binance "github.com/laoliu6668/esharp_binance_utils"
)

// # MODEL 获取用户账户
type ApiResponseAccountData struct {
	Data []AccountData `json:"balances"`
}

type AccountData struct {
	Asset  string `json:"asset"`
	Free   string `json:"free"`
	Locked string `json:"locked"`
}

// 获取现货账户信息
// doc: https://binance-docs.github.io/apidocs/spot/cn/#user_data-42
func GetSpotAccount() (data []AccountData, err error) {
	const flag = "binance GetSpotAccount"
	body, _, err := binance.ApiConfig.Get(gateway_binance, "/api/v3/account", nil)
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}
	// util.WriteTestJsonFile(flag, body)
	res := ApiResponseAccountData{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}
	return res.Data, nil
}
