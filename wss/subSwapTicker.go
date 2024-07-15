package binance_wss

import (
	"encoding/json"
	"fmt"
	"strings"

	root "github.com/laoliu6668/esharp_binance_utils"
	"github.com/laoliu6668/esharp_binance_utils/util"
	"github.com/laoliu6668/esharp_binance_utils/util/websocketclient"
)

type WssSwapMsg struct {
	E  string `json:"e"` // 事件名
	U  int64  `json:"u"` // 更新ID
	S  string `json:"s"` // 交易对
	B  string `json:"b"` // 买一价
	BS string `json:"B"` // 买一量
	A  string `json:"a"` // 卖一价
	AS string `json:"A"` // 卖一量
	T  int64  `json:"t"` // 交易时间
	ET int64  `json:"E"` // 事件时间 0.0000000848
}

func SubSwapTicker(reciveHandle func(Ticker), logHandle func(string), errHandle func(error)) {
	gateway := "wss://fstream.binance.com/ws/!bookTicker"
	proxyUrl := ""
	if root.UseProxy {
		proxyUrl = fmt.Sprintf("http://%s", root.ProxyUrl)
	}
	ws := websocketclient.New(gateway, proxyUrl)
	ws.OnConnectError(func(err error) {
		go errHandle(fmt.Errorf("OnConnectError: %v", err))
	})
	ws.OnDisconnected(func(err error) {
		go errHandle(fmt.Errorf("disconnected: %v", err))
	})
	ws.OnConnected(func() {
		go logHandle("SubSwapTicker Connected")
	})
	ws.OnTextMessageReceived(func(msg string) {
		m := WssSwapMsg{}
		err := json.Unmarshal([]byte(msg), &m)
		if err != nil {
			go errHandle(fmt.Errorf("msg json.Unmarshal err: %s", msg))
			return
		}
		if m.E == "bookTicker" {
			go reciveHandle(Ticker{
				Exchange: root.ExchangeName,
				Symbol:   strings.Replace(m.S, "USDT", "", 1),
				Buy: Values{
					Price: util.ParseFloat(m.B, 0),
					Size:  util.ParseFloat(m.BS, 0),
				},
				Sell: Values{
					Price: util.ParseFloat(m.A, 0),
					Size:  util.ParseFloat(m.AS, 0),
				},
				UpdateAt: root.GetTimeFloat(),
			})
		} else {
			go logHandle("unkown msg: " + msg)
		}
	})

	ws.OnClose(func(code int, text string) {
		go errHandle(fmt.Errorf("close: %v, %v", code, text))
	})

	ws.Connect()

}
