package apis

import (
	"encoding/json"
	"fmt"
	"strings"

	root "github.com/laoliu6668/esharp_binance_utils"
	"github.com/laoliu6668/esharp_binance_utils/util"
)

type Res struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

// 变换逐全仓模式 (TRADE)
// param margin_type: ISOLATED(逐仓), CROSSED(全仓)
// doc:  https://developers.binance.com/docs/zh-CN/derivatives/usds-margined-futures/trade/rest-api/Change-Margin-Type
func ChangeSwapMarginType(symbol string, ISOLATED bool) (err error) {
	const flag = "binance ChangeSwapMarginType"
	marginType := "CROSSED"
	if ISOLATED {
		marginType = "ISOLATED"
	}
	body, _, err := root.ApiConfig.Post(gateway_fapi, "/fapi/v1/marginType", map[string]any{
		"marginType": marginType,
		"symbol":     fmt.Sprintf("%sUSDT", strings.ToUpper(symbol)),
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		return
	}
	// fmt.Printf("body: %v\n", body)

	res := Res{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}

	if res.Code != 200 {
		err = fmt.Errorf("%s err: %v", flag, res.Msg)
		return
	}

	return nil
}

// 更改持仓模式(TRADE)
// param dual_side_position: "true": 双向持仓模式；"false": 单向持仓模式
// doc: https://developers.binance.com/docs/zh-CN/derivatives/usds-margined-futures/trade/rest-api/Change-Position-Mode
func ChangeSwapPositionSideDual(dual_side_position bool) (err error) {
	const flag = "binance ChangeWwapPositionSideDual"
	body, _, err := root.ApiConfig.Post(gateway_fapi, "/fapi/v1/positionSide/dual", map[string]any{
		"dualSidePosition": dual_side_position,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		return
	}

	res := Res{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}

	if res.Code != 200 {
		err = fmt.Errorf("%s err: %v", flag, res.Msg)
		return
	}

	return nil
}

// 调整开仓杠杆(TRADE)
// param leverage: 目标杠杆倍数：1 到 125 整数
// doc: https://developers.binance.com/docs/zh-CN/derivatives/usds-margined-futures/trade/rest-api/Change-Initial-Leverage
func ChangeSwapLeverage(symbol string, leverage int) (err error) {
	const flag = "binance ChangeSwapLeverage"
	body, _, err := root.ApiConfig.Post(gateway_fapi, "/fapi/v1/leverage", map[string]any{
		"symbol":   fmt.Sprintf("%sUSDT", strings.ToUpper(symbol)),
		"leverage": leverage,
	})
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		return
	}

	res := Res{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}

	if res.Code != 200 {
		err = fmt.Errorf("%s err: %v", flag, res.Msg)
		return
	}

	return nil
}

type SwapBalance struct {
	AccountAlias       string `json:"accountAlias"`
	Asset              string `json:"asset"`
	Balance            string `json:"balance"`
	CrossWalletBalance string `json:"crossWalletBalance"`
	CrossUnPnl         string `json:"crossUnPnl"`
	AvailableBalance   string `json:"availableBalance"`
	MaxWithdrawAmount  string `json:"maxWithdrawAmount"`
	MarginAvailable    bool   `json:"marginAvailable"`
	UpdateTime         int    `json:"updateTime"`
}

// 账户余额V2
// doc: https://binance-docs.github.io/apidocs/futures/cn/#v2-user_data
// symbol default:USDT
func GetSwapBalance(symbol string) (b SwapBalance, err error) {
	const flag = "binance GetSwapBalance"
	body, _, err := root.ApiConfig.Get(gateway_fapi, "/fapi/v2/balance", nil)
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		return
	}
	// fmt.Printf("body: %s\n", body)
	res := []SwapBalance{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}
	if symbol == "" {
		symbol = "USDT"
	}
	has := false
	for _, v := range res {
		if v.Asset == symbol {
			has = true
			b = v
		}
	}
	if !has {
		err = fmt.Errorf("%s err: %v", flag, "symbol not found")
		return
	}

	return
}

type SwapPosition struct {
	PositionAmt            string `json:"positionAmt"`    // 持仓数量
	PositionSide           string `json:"positionSide"`   // 持仓方向
	BreakEvenPrice         string `json:"breakEvenPrice"` // 持仓成本价
	Leverage               string `json:"leverage"`       // 杠杆倍率
	Isolated               bool   `json:"isolated"`       // 是否是逐仓模式
	EntryPrice             string `json:"entryPrice"`
	InitialMargin          string `json:"initialMargin"`
	IsolatedWallet         string `json:"isolatedWallet"`
	MaintMargin            string `json:"maintMargin"`
	MaxNotional            string `json:"maxNotional"`
	Notional               string `json:"notional"`
	OpenOrderInitialMargin string `json:"openOrderInitialMargin"`
	PositionInitialMargin  string `json:"positionInitialMargin"`
	Symbol                 string `json:"symbol"`
	UnrealizedProfit       string `json:"unrealizedProfit"`
	UpdateTime             int64  `json:"updateTime"`
}
type SwapAccount struct {
	AvailableBalance            string         `json:"availableBalance"`        // 以USD计价的可用余额
	TotalCrossWalletBalance     string         `json:"totalCrossWalletBalance"` // 以USD计价的全仓账户余额
	TotalInitialMargin          string         `json:"totalInitialMargin"`      // 以USD计价的所需起始保证金总额
	TotalMaintMargin            string         `json:"totalMaintMargin"`        // 以USD计价的维持保证金总额
	TotalMarginBalance          string         `json:"totalMarginBalance"`      // 以USD计价的保证金总余额
	TotalWalletBalance          string         `json:"totalWalletBalance"`      // 以USD计价的账户总余额
	CanDeposit                  bool           `json:"canDeposit"`              // 是否可以入金
	CanTrade                    bool           `json:"canTrade"`                // 是否可以交易
	CanWithdraw                 bool           `json:"canWithdraw"`             // 是否可以出金
	FeeBurn                     bool           `json:"feeBurn"`                 // "true": 手续费抵扣开; "false": 手续费抵扣关
	FeeTier                     int            `json:"feeTier"`
	MaxWithdrawAmount           string         `json:"maxWithdrawAmount"`
	MultiAssetsMargin           bool           `json:"multiAssetsMargin"`
	TotalCrossUnPnl             string         `json:"totalCrossUnPnl"`
	TotalOpenOrderInitialMargin string         `json:"totalOpenOrderInitialMargin"`
	TotalPositionInitialMargin  string         `json:"totalPositionInitialMargin"`
	TotalUnrealizedProfit       string         `json:"totalUnrealizedProfit"`
	TradeGroupID                int            `json:"tradeGroupId"`
	UpdateTime                  int            `json:"updateTime"`
	Posistions                  []SwapPosition `json:"positions"` // 持仓
	// Assets                      []SwapPosition `json:"assets"` // 资产
}

// 账户信息 持仓
// doc: https://binance-docs.github.io/apidocs/futures/cn/#v2-user_data-2
func GetSwapAccount() (data SwapAccount, err error) {
	const flag = "binance GetSwapBalance"
	body, _, err := root.ApiConfig.Get(gateway_fapi, "/fapi/v2/account", nil)
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		return
	}
	// fmt.Printf("body: %s\n", body)
	err = json.Unmarshal(body, &data)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}
	return
}

type SwapFunding struct {
	Symbol      string `json:"symbol"`          // 交易对 "BTCUSDT"
	FundingRate string `json:"lastFundingRate"` // 最近更新的资金费率
	// "nextFundingTime": 1597392000000,   // 下次资金费时间
	// "interestRate": "0.00010000",       // 标的资产基础利率
	Time int64 `json:"time"` // 更新时间 1597370495002
}

// 期货资金费率
// doc: https://developers.binance.com/docs/zh-CN/derivatives/usds-margined-futures/market-data/rest-api/Get-Funding-Info
func GetSwapFunding() (data []SwapFunding, err error) {
	const flag = "binance GetSwapFunding"
	body, _, err := root.ApiConfig.Request("GET", gateway_fapi, "/fapi/v1/premiumIndex", nil, 0, false)
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		return
	}
	// fmt.Printf("body: %s\n", body)
	err = json.Unmarshal(body, &data)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}
	return
}

// 账户信息 持仓风险
// doc: https://binance-docs.github.io/apidocs/futures/cn/#v2-user_data-2
func GetPositionRisk() (data SwapAccount, err error) {
	const flag = "binance GetPositionRisk"
	body, _, err := root.ApiConfig.Get(gateway_fapi, "/fapi/v2/positionRisk", nil)
	if err != nil {
		err = fmt.Errorf("%s err: %v", flag, err)
		return
	}
	fmt.Printf("body: %s\n", body)
	util.WriteTestJsonFile(flag, body)
	err = json.Unmarshal(body, &data)
	if err != nil {
		err = fmt.Errorf("%s jsonDecodeErr: %v", flag, err)
		fmt.Println(err)
		return
	}
	return
}
