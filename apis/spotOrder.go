package apis

import (
	"encoding/json"
	"fmt"

	root "github.com/laoliu6668/esharp_binance_utils"
)

type ApiResponseListData struct {
	Status  string           `json:"status"`
	Message string           `json:"err_msg"`
	Data    []map[string]any `json:"data"`
}

func (a *ApiResponseListData) Success() bool {
	return a.Status == "ok"
}

// ### 现货下单
// doc: https://developers.binance.com/docs/zh-CN/binance-spot-api-docs/rest-api#%E4%B8%8B%E5%8D%95-trade

type OrderRes struct {
	OrderID int `json:"orderId"`
}

func SpotBuyMarket(symb string, amount float64) (orderId string, err error) {
	// 市价买入
	const flag = "Binance SpotBuyMarket"
	body, _, err := root.ApiConfig.Post(gateway_binance, "/api/v3/order", map[string]any{
		"symbol":           symb + "USDT",
		"side":             "BUY",
		"type":             "MARKET",
		"quantity":         amount,
		"newOrderRespType": "ACK",
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}
	// fmt.Printf("string(body): %v\n", string(body))
	res := OrderRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}

	return fmt.Sprintf("%v", res.OrderID), nil
}

func SpotSellMarket(symb string, volume float64) (data string, err error) {
	// 市价卖出
	const flag = "Binance SpotSellMarket"
	body, _, err := root.ApiConfig.Post(gateway_binance, "/api/v3/order", map[string]any{
		"symbol":           symb + "USDT",
		"side":             "SELL",
		"type":             "MARKET",
		"quantity":         volume,
		"newOrderRespType": "ACK",
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		fmt.Println(err)
		return
	}
	// fmt.Printf("string(body): %v\n", string(body))
	res := OrderRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}

	return fmt.Sprintf("%v", res.OrderID), nil
}
