package apis

import (
	"encoding/json"
	"fmt"

	binance "github.com/laoliu6668/esharp_binance_utils"
)

const gateway_binance = "api.binance.com"

type TranId struct {
	TranId int `json:"tranId"`
}

// ### 现货账户向期货账户划转
// doc: https://binance-docs.github.io/apidocs/spot/cn/#user_data-14
func SpotToSwapTransfer(amount float64) (id string, err error) {
	const flag = "binance SpotToSwapTransfer"
	body, _, err := binance.ApiConfig.Post(gateway_binance, "/sapi/v1/asset/transfer", map[string]any{
		"type":   "MAIN_UMFUTURE",
		"asset":  "USDT",
		"amount": amount,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}

	res := TranId{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}

	return fmt.Sprintf("%v", res.TranId), nil
}

// 期货账户向现货账户划转
// doc: https://binance-docs.github.io/apidocs/spot/cn/#user_data-14
func SwapToSpotTransfer(amount float64) (id string, err error) {
	const flag = "binance SwapToSpotTransfer"
	body, _, err := binance.ApiConfig.Post(gateway_binance, "/sapi/v1/asset/transfer", map[string]any{
		"type":   "UMFUTURE_MAIN",
		"asset":  "USDT",
		"amount": amount,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		return
	}
	// fmt.Printf("body: %s\n", body)
	res := TranId{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}

	return fmt.Sprintf("%v", res.TranId), nil

}
