package apis

import (
	"encoding/json"
	"fmt"
	"strings"

	root "github.com/laoliu6668/esharp_binance_utils"
)

const gateway_fapi = "fapi.binance.com"

// 期货卖出开空
// doc: https://developers.binance.com/docs/zh-CN/derivatives/usds-margined-futures/trade/rest-api/New-Order
func SwapSellOpen(symb string, volume int) (orderId string, err error) {
	const symbol = "HTX SwapSellOpen"
	body, _, err := root.ApiConfig.Post(gateway_fapi, "/fapi/v1/order", map[string]any{
		"symbol":       fmt.Sprintf("%sUSDT", strings.ToUpper(symb)),
		"volume":       volume,
		"side":         "SELL",
		"positionSide": "SHORT",
		"type":         "MARKET",
		"quantity":     volume,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	// fmt.Printf("string(body): %v\n", string(body))
	res := OrderRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}

	return fmt.Sprintf("%v", res.OrderID), nil
}

// 期货买入平空
// doc: https://developers.binance.com/docs/zh-CN/derivatives/usds-margined-futures/trade/rest-api/New-Order
func SwapBuyClose(symb string, volume int) (orderId string, err error) {
	const symbol = "Binance SwapBuyClose"
	body, _, err := root.ApiConfig.Post(gateway_fapi, "/fapi/v1/order", map[string]any{
		"symbol":       fmt.Sprintf("%sUSDT", strings.ToUpper(symb)),
		"volume":       volume,
		"side":         "BUY",
		"positionSide": "SHORT",
		"type":         "MARKET",
		"quantity":     volume,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	// fmt.Printf("string(body): %v\n", string(body))
	res := OrderRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	return fmt.Sprintf("%v", res.OrderID), nil
}

// 期货买入开多
// doc: https://developers.binance.com/docs/zh-CN/derivatives/usds-margined-futures/trade/rest-api/New-Order
func SwapBuyOpen(symb string, volume int) (orderId string, err error) {
	const symbol = "Binance SwapBuyOpen"
	body, _, err := root.ApiConfig.Post(gateway_fapi, "/fapi/v1/order", map[string]any{
		"symbol":       fmt.Sprintf("%sUSDT", strings.ToUpper(symb)),
		"volume":       volume,
		"side":         "BUY",
		"positionSide": "LONG",
		"type":         "MARKET",
		"quantity":     volume,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	// fmt.Printf("string(body): %v\n", string(body))
	res := OrderRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	return fmt.Sprintf("%v", res.OrderID), nil
}

// 期货买入开多
// doc: https://developers.binance.com/docs/zh-CN/derivatives/usds-margined-futures/trade/rest-api/New-Order
func SwapSellClose(symb string, volume int) (orderId string, err error) {
	const symbol = "Binance SwapSellClose"
	body, _, err := root.ApiConfig.Post(gateway_fapi, "/fapi/v1/order", map[string]any{
		"symbol":       fmt.Sprintf("%sUSDT", strings.ToUpper(symb)),
		"volume":       volume,
		"side":         "SELL",
		"positionSide": "LONG",
		"type":         "MARKET",
		"quantity":     volume,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	// fmt.Printf("string(body): %v\n", string(body))
	res := OrderRes{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	return fmt.Sprintf("%v", res.OrderID), nil
}

// 增加空头逐仓保证金
// doc: https://binance-docs.github.io/apidocs/futures/cn/#trade-14
func SwapIncShortPositionMargin(symb string, amount int) (err error) {
	const symbol = "Binance SwapIncShortPositionMargin"
	body, _, err := root.ApiConfig.Post(gateway_fapi, "/fapi/v1/positionMargin", map[string]any{
		"symbol":       fmt.Sprintf("%sUSDT", strings.ToUpper(symb)),
		"amount":       amount,
		"positionSide": "SHORT",
		"type":         1, // type 1: 增加逐仓保证金，2: 减少逐仓保证金
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	// fmt.Printf("string(body): %v\n", string(body))
	res := Res{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if res.Code != 200 {
		return fmt.Errorf("%v", res.Msg)
	}
	return nil
}

// 减少空头逐仓保证金
// doc: https://binance-docs.github.io/apidocs/futures/cn/#trade-14
func SwapDecShortPositionMargin(symb string, amount int) (err error) {
	const symbol = "Binance SwapDecShortPositionMargin"
	body, _, err := root.ApiConfig.Post(gateway_fapi, "/fapi/v1/positionMargin", map[string]any{
		"symbol":       fmt.Sprintf("%sUSDT", strings.ToUpper(symb)),
		"amount":       amount,
		"positionSide": "SHORT",
		"type":         2, // type 1: 增加逐仓保证金，2: 减少逐仓保证金
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	res := Res{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if res.Code != 200 {
		return fmt.Errorf("%v", res.Msg)
	}
	return nil
}

// 增加多头逐仓保证金
// doc: https://binance-docs.github.io/apidocs/futures/cn/#trade-14
func SwapIncLongPositionMargin(symb string, amount int) (err error) {
	const symbol = "Binance SwapIncLongPositionMargin"
	body, _, err := root.ApiConfig.Post(gateway_fapi, "/fapi/v1/positionMargin", map[string]any{
		"symbol":       fmt.Sprintf("%sUSDT", strings.ToUpper(symb)),
		"amount":       amount,
		"positionSide": "LONG",
		"type":         1, // type 1: 增加逐仓保证金，2: 减少逐仓保证金
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	// fmt.Printf("string(body): %v\n", string(body))
	res := Res{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if res.Code != 200 {
		return fmt.Errorf("%v", res.Msg)
	}
	return nil
}

// 减少多头逐仓保证金
// doc: https://binance-docs.github.io/apidocs/futures/cn/#trade-14
func SwapDecLongPositionMargin(symb string, amount int) (err error) {
	const symbol = "Binance SwapDecLongPositionMargin"
	body, _, err := root.ApiConfig.Post(gateway_fapi, "/fapi/v1/positionMargin", map[string]any{
		"symbol":       fmt.Sprintf("%sUSDT", strings.ToUpper(symb)),
		"amount":       amount,
		"positionSide": "LONG",
		"type":         2, // type 1: 增加逐仓保证金，2: 减少逐仓保证金
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	fmt.Printf("string(body): %v\n", string(body))
	res := Res{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	if res.Code != 200 {
		return fmt.Errorf("%v", res.Msg)
	}
	return nil
}

// 撤销全部订单
func CancelALLOrder(symb string) (err error) {
	const symbol = "HTX CancelALLOrder"
	body, _, err := root.ApiConfig.Delete(gateway_fapi, "/fapi/v1/allOpenOrders", map[string]any{
		"symbol": fmt.Sprintf("%sUSDT", strings.ToUpper(symb)),
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", symbol, err)
		fmt.Println(err)
		return
	}
	fmt.Printf("string(body): %v\n", string(body))
	res := Res{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", symbol, err)
		fmt.Println(err)
		return
	}
	return nil
}
