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

const proto = "wss://"

type balanceRes struct {
	Symbol string `json:"a"` // 资产名称
	Free   string `json:"f"` // 可用余额
	Lokc   string `json:"l"` // 冻结余额
}
type order struct {
	Symbol      string `json:"s"`
	OrderFaq    string `json:"S"` // 订单方向 BUY SELL
	OrderId     int64  `json:"i"` // orderId
	OrderStatus string `json:"X"` // orderStatus
	OrderVolume string `json:"q"` // 订单原始数量
	OrderPrice  string `json:"p"` // 订单原始价格
	TradeVolume string `json:"z"` // 订单累计已成交量
	TradeValue  string `json:"Z"` // 订单累计已成交金额
	OrderType   string `json:"o"` // 订单类型
	CreatedAt   int64  `json:"O"` // 订单创建时间
	FilledAt    int64  `json:"T"` // 订单成交时间
}
type data struct {
	Event     string `json:"e"` // 事件类型
	EventTime int64  `json:"E"` // 事件时间
}
type account struct {
	Event     string       `json:"e"` // 事件类型
	EventTime int64        `json:"E"` // 事件时间
	Symbol    string       `json:"s"` // 账户本次更新时间戳
	B         []balanceRes `json:"B"` // 余额 outboundAccountPosition
}
type msgOrder struct {
	Stream string `json:"stream"`
	Data   order  `json:"data"`
}
type msgAccount struct {
	Stream string  `json:"stream"`
	Data   account `json:"data"`
}
type msg struct {
	Stream string `json:"stream"`
	Data   data   `json:"data"`
}

// 订阅现货账户变化
// 同步:ReciveBalanceMsg 异步:reciveOrderHandle、logHandle、errHandle
func SubSpotAccount(reciveAccHandle func(ReciveBalanceMsg), reciveOrderHandle func(ReciveSpotOrderMsg), logHandle func(string), errHandle func(error)) {
	const flag = "Binance SubAccountUpdate"
	s, err := apis.GetSpotAccountListenKey()
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}
	fmt.Printf("listenKey: %v\n", s)
	gateway := "stream.binance.com:9443/stream?streams="

	// 保持listenKey

	ticker := time.NewTicker(time.Minute * 30)

	// defer ticker.Stop()
	go func() {
		for range ticker.C {
			fmt.Printf("keep listenKey: %v\n", s)
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
		m := msg{}
		err := json.Unmarshal([]byte(message), &m)
		if err != nil {
			go errHandle(fmt.Errorf("%s json.Unmarshal %s %s", flag, err.Error(), message))
			return
		}
		if m.Data.Event == "outboundAccountPosition" {
			m := msgAccount{}
			err := json.Unmarshal([]byte(message), &m)
			if err != nil {
				go errHandle(fmt.Errorf("%s json.Unmarshal %s %s", flag, err.Error(), message))
				return
			}
			for _, v := range m.Data.B {
				m := msgAccount{}
				err := json.Unmarshal([]byte(message), &m)
				if err != nil {
					go errHandle(fmt.Errorf("%s json.Unmarshal %s %s", flag, err.Error(), message))
					return
				}
				// 账户更新
				reciveAccHandle(ReciveBalanceMsg{
					Exchange: root.ExchangeName,
					Symbol:   v.Symbol,
					Free:     util.ParseFloat(v.Free, 0),
					Lock:     util.ParseFloat(v.Lokc, 0),
				})
			}

		} else if m.Data.Event == "executionReport" {
			m := msgOrder{}
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
			tradeValueD, _ := decimal.NewFromString(o.TradeValue)
			tradeVolumeD, _ := decimal.NewFromString(o.TradeVolume)
			tradePriceD := tradeValueD.Div(tradeVolumeD)
			tradePrice, _ := tradePriceD.Float64()
			// 订单更新
			go reciveOrderHandle(ReciveSpotOrderMsg{
				Exchange:    root.ExchangeName,
				Symbol:      symbol,
				OrderId:     "B" + strconv.FormatInt(o.OrderId, 10),
				OrderType:   strings.ToLower(o.OrderFaq + "-" + o.OrderType),
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
