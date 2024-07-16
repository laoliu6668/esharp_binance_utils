package wss

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	root "github.com/laoliu6668/esharp_binance_utils"
	"github.com/laoliu6668/esharp_binance_utils/apis"
	"github.com/laoliu6668/esharp_binance_utils/util"
	"github.com/laoliu6668/esharp_binance_utils/util/websocketclient"
	"github.com/shopspring/decimal"
)

const proto = "wss://"

type balanceRes struct {
	Symbol string `json:"a"` // 资产名称
	Free   string `json:"f"` // 可用余额
	Lokc   string `json:"l"` // 冻结余额
}

// 订阅现货账户变化
func SubSpotAccount(reciveAccHandle func(ReciveBalanceMsg), reciveOrderHandle func(ReciveSpotOrderMsg), logHandle func(string), errHandle func(error)) {
	const flag = "Binance SubAccountUpdate"
	s, err := apis.GetSpotAccountListenKey()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("s: %v\n", s)
	gateway := "stream.binance.com:9443/stream?streams="

	// 保持listenKey

	ticker := time.NewTicker(time.Minute * 30)

	defer ticker.Stop()
	go func() {
		for range ticker.C {
			apis.KeepSpotAccountListenKey(s)
		}
	}()

	requrl := fmt.Sprintf("%s%s%s", proto, gateway, s)
	proxyUrl := ""
	if root.UseProxy {
		go logHandle(fmt.Sprintf("proxyUrl: %v\n", root.ProxyUrl))
		proxyUrl = fmt.Sprintf("http://%s", root.ProxyUrl)

	}
	ws := websocketclient.New(requrl, proxyUrl)
	ws.OnConnectError(func(err error) {
		go errHandle(err)
	})
	ws.OnDisconnected(func(err error) {
		go errHandle(err)
	})
	ws.OnConnected(func() {
		go logHandle(fmt.Sprintf("%s connected", flag))
	})
	ws.OnTextMessageReceived(func(message string) {
		fmt.Printf("message: %v\n", message)

		type data struct {
			Event       string     `json:"e"` // 事件类型
			EventTime   int64      `json:"E"` // 事件时间
			Symbol      string     `json:"s"` // 账户本次更新时间戳
			B           balanceRes `json:"B"` // 余额 outboundAccountPosition
			OrderFaq    string     `json:"S"` // 订单方向 BUY SELL
			OrderId     string     `json:"i"` // orderId
			OrderStatus string     `json:"X"` // orderStatus
			OrderVolume string     `json:"q"` // 订单原始数量
			OrderPrice  string     `json:"p"` // 订单原始价格
			TradeVolume string     `json:"z"` // 订单累计已成交量
			TradeValue  string     `json:"Z"` // 订单累计已成交金额
			CreatedAt   int64      `json:"O"` // 订单创建时间
			FilledAt    int64      `json:"T"` // 订单成交时间

		}
		type msg struct {
			Stream string `json:"stream"`
			Data   data   `json:"data"`
		}
		m := msg{}
		err := json.Unmarshal([]byte(message), &m)
		if err != nil {
			go errHandle(fmt.Errorf("%s json.Unmarshal %s", flag, err.Error()))
			return
		}
		symbol := strings.Replace(m.Data.Symbol, "USDT", "", 1)
		if m.Data.Event == "outboundAccountPosition" {
			// 账户更新
			reciveAccHandle(ReciveBalanceMsg{
				Exchange: root.ExchangeName,
				Symbol:   symbol,
				Free:     util.ParseFloat(m.Data.B.Free, 0),
				Lock:     util.ParseFloat(m.Data.B.Lokc, 0),
			})
		} else if m.Data.Event == "executionReport" {
			if m.Data.OrderStatus != "FILLED" {
				// 只要成交订单
				return
			}
			o := m.Data
			orderPriceD, _ := decimal.NewFromString(o.OrderPrice)
			orderVolumeD, _ := decimal.NewFromString(o.OrderVolume)
			orderValueD := orderPriceD.Mul(orderVolumeD)
			orderValue, _ := orderValueD.Float64()
			tradeValueD, _ := decimal.NewFromString(o.TradeValue)
			tradeVolumeD, _ := decimal.NewFromString(o.TradeVolume)
			tradePriceD := tradeValueD.Div(tradeVolumeD)
			tradePrice, _ := tradePriceD.Float64()
			// 订单更新
			reciveOrderHandle(ReciveSpotOrderMsg{
				Exchange:    root.ExchangeName,
				Symbol:      symbol,
				OrderId:     o.OrderId,
				OrderType:   strings.ToLower(o.OrderFaq) + "-market",
				OrderPrice:  util.ParseFloat(o.OrderPrice, 0),
				OrderVolume: util.ParseFloat(o.OrderVolume, 0),
				OrderValue:  orderValue,
				TradePrice:  tradePrice,
				TradeVolume: util.ParseFloat(o.TradeVolume, 0),
				TradeValue:  util.ParseFloat(o.TradeValue, 0),
				CreateAt:    o.CreatedAt,
				FilledAt:    o.FilledAt,
				Status:      2,
			})
		} else if m.Data.Event == "balanceUpdate" {
			logHandle(message)
		}
	})

	ws.OnClose(func(code int, text string) {
		go errHandle(fmt.Errorf("close: %v, %v", code, text))
	})

	ws.Connect()

}
