package wss

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	root "github.com/laoliu6668/esharp_binance_utils"
	"github.com/laoliu6668/esharp_binance_utils/apis"
	"github.com/laoliu6668/esharp_binance_utils/util"
	"github.com/laoliu6668/esharp_binance_utils/util/websocketclient"
	"github.com/shopspring/decimal"
)

type SwapBalance struct {
	Asset            string `json:"a"` // 资产名称 USDT
	Balance          string `json:"wb"`
	AvailableBalance string `json:"cw"`
}
type SwapPosition struct {
	Symbol string `json:"s"`  // 交易对 BTCUSDT
	Faq    string `json:"ps"` // 持仓方向 LONG SHORT
	Volume string `json:"pa"` // 仓位 -10
	Margin string `json:"iw"` // 持仓保证金
}
type SwapAccount struct {
	Symbol    string         `json:"s"`
	Balances  []SwapBalance  `json:"B"`
	Positions []SwapPosition `json:"P"`
}

type SwapMargin struct {
	Symbol      string  `json:"s"`
	Faq         string  `json:"ps"`           // 持仓方向 LONG SHORT
	Volume      string  `json:"pa"`           // 仓位
	Margin      string  `json:"iw"`           // 若为逐仓，仓位保证金
	KeepMNargin string  `json:"mm"`           // 持仓需要的维持保证金
	MarginRatio float64 `json:"margin_ratio"` // 保证金率
}
type msgSwapMargin struct {
	Margins []SwapMargin `json:"p"`
}
type msgSwapAccount struct {
	Data SwapAccount `json:"a"`
}
type swapOrder struct {
	Symbol           string `json:"s"`
	OrderFaq         string `json:"S"`  // 订单方向 BUY SELL
	OrderId          int64  `json:"i"`  // orderId
	OrderStatus      string `json:"X"`  // 订单状态 FILLED
	OrderVolume      string `json:"q"`  // 订单原始数量
	OrderVolumePlace string `json:"Q"`  // 订单原始数量Place
	OrderPrice       string `json:"p"`  // 订单原始价格
	OrderPricePlace  string `json:"P"`  // 订单原始价格Place
	TradeVolume      string `json:"z"`  // 订单累计已成交量
	TradePrice       string `json:"ap"` // 订单平均价格
	Place            string `json:"Z"`  // 订单累计已成交量place
	OrderType        string `json:"o"`  // 订单类型 MARKET
	CreatedAt        int64  `json:"O"`  // 订单创建时间
	FilledAt         int64  `json:"T"`  // 订单成交时间
	TID              int64  `json:"t"`  // 订单成交id
	Faq              string `json:"ps"` // 持仓方向 LONG SHORT
}
type msgSwapOrder struct {
	Data swapOrder `json:"o"`
}

// 订阅期货账户变化
// 余额+持仓+订单
func SubSwapAccount(reciveAccHandle func(SwapAccount), reciveMarginHandle func(SwapMargin), reciveOrderHandle func(ReciveSwapOrderMsg), logHandle func(string), errHandle func(error)) {
	const flag = "SubSwapAccount"
	var listenKey string
	for {
		s, err := apis.GetSwapAccountListenKey()
		if err != nil {
			errHandle(fmt.Errorf("GetSpotAccountListenKey: %v", err.Error()))
			time.Sleep(time.Second * 10)
			continue
		}
		listenKey = s
		break
	}
	ticker := time.NewTicker(time.Minute * 30)
	go func() {
		for range ticker.C {
			apis.KeepSwapAccountListenKey(listenKey)
		}
	}()

	gateway := "fstream.binance.com?/ws/"

	requrl := fmt.Sprintf("wss://%s%s", gateway, listenKey)
	proxyUrl := ""
	if root.UseProxy {
		go logHandle(fmt.Sprintf("proxyUrl: %v\n", root.ProxyUrl))
		proxyUrl = fmt.Sprintf("http://%s", root.ProxyUrl)
	}
	go logHandle(fmt.Sprintf("requrl: %v\n", requrl))
	ws := websocketclient.New(requrl, proxyUrl)
	ws.OnConnectError(func(err error) {
		go errHandle(err)
	})
	ws.OnDisconnected(func(err error) {
		go errHandle(err)
	})
	ws.OnConnected(func() {
		go logHandle("connected Socket")
	})
	ws.OnTextMessageReceived(func(message string) {
		m := data{}
		err := json.Unmarshal([]byte(message), &m)
		if err != nil {
			go errHandle(fmt.Errorf("%s json.Unmarshal %s %s", flag, err.Error(), message))
			return
		}
		if m.Event == "ACCOUNT_UPDATE" {
			m := msgSwapAccount{}
			err := json.Unmarshal([]byte(message), &m)
			if err != nil {
				go errHandle(fmt.Errorf("%s json.Unmarshal %s %s", flag, err.Error(), message))
				return
			}
			go reciveAccHandle(m.Data)
		} else if m.Event == "ORDER_TRADE_UPDATE" {
			fmt.Printf("message: %v\n", message)
			m := msgSwapOrder{}
			err := json.Unmarshal([]byte(message), &m)
			if err != nil {
				go errHandle(fmt.Errorf("%s json.Unmarshal %s %s", flag, err.Error(), message))
				return
			}
			if m.Data.OrderStatus != "FILLED" {
				// 只要成交订单
				return
			}
			symbol := strings.Replace(m.Data.Symbol, "USDT", "", 1)
			if m.Data.Symbol == "USDT" {
				symbol = "USDT"
			}
			o := m.Data
			orderPriceD, _ := decimal.NewFromString(o.OrderPrice)
			orderVolumeD, _ := decimal.NewFromString(o.OrderVolume)
			orderValueD := orderPriceD.Mul(orderVolumeD)
			orderValue, _ := orderValueD.Float64()
			tradePriceD, _ := decimal.NewFromString(o.TradePrice)
			tradePrice, _ := tradePriceD.Float64()
			tradeVolumeD, _ := decimal.NewFromString(o.TradeVolume)
			tradeVolume, _ := tradeVolumeD.Float64()
			tradeValueD := tradePriceD.Mul(tradeVolumeD) // 成交金额 = 成交价 * 成交量
			tradeValue, _ := tradeValueD.Float64()
			// 订单更新
			orderType := ""
			if o.OrderFaq == "BUY" && o.Faq == "LONG" {
				orderType = "buy-open"
			} else if o.OrderFaq == "BUY" && o.Faq == "SHORT" {
				orderType = "buy-close"
			} else if o.OrderFaq == "SELL" && o.Faq == "LONG" {
				orderType = "sell-close"
			} else if o.OrderFaq == "SELL" && o.Faq == "SHORT" {
				orderType = "sell-open"
			}
			go reciveOrderHandle(ReciveSwapOrderMsg{
				Exchange:    root.ExchangeName,
				Symbol:      symbol,
				OrderId:     strconv.FormatInt(o.OrderId, 10),
				OrderType:   orderType,
				OrderPrice:  util.ParseFloat(o.OrderPrice, 0),
				OrderVolume: util.ParseFloat(o.OrderVolume, 0),
				OrderValue:  orderValue,
				TradePrice:  tradePrice,
				TradeVolume: tradeVolume,
				TradeValue:  tradeValue,
				CreateAt:    o.CreatedAt,
				FilledAt:    o.FilledAt,
				Status:      2,
			})
		} else if m.Event == "MARGIN_CALL" {
			m := msgSwapMargin{}
			err := json.Unmarshal([]byte(message), &m)
			if err != nil {
				go errHandle(fmt.Errorf("%s json.Unmarshal %s %s", flag, err.Error(), message))
				return
			}

			for _, v := range m.Margins {
				margin := util.ParseFloat(v.Margin, 0)
				KeepMNargin := util.ParseFloat(v.KeepMNargin, 1)
				v.MarginRatio = util.FixedFloat(margin/KeepMNargin, 0)
				go reciveMarginHandle(v)
			}
		} else {
			go logHandle(message)
		}
	})
	ws.OnClose(func(code int, text string) {
		// fmt.Printf("close: %v, %v\n", code, text)
		go errHandle(fmt.Errorf("close: %v, %v", code, text))
	})
	ws.Connect()

}
