package wss

import (
	"encoding/json"
	"fmt"
	"strings"

	root "github.com/laoliu6668/esharp_binance_utils"
	"github.com/laoliu6668/esharp_binance_utils/util"
	"github.com/laoliu6668/esharp_binance_utils/util/websocketclient"
)

type WssSpotMsg struct {
	Id        string `json:"id"`
	Symbol    string `json:"s"` // 交易对
	BuyPrice  string `json:"b"` // 买单最优挂单价格
	BuySize   string `json:"B"` // 买单最优挂单数量
	SellPrice string `json:"a"` // 卖单最优挂单价格
	SellSize  string `json:"A"` // 卖单最优挂单数量
}

func SubSpotTicker(symbols []string, reciveHandle func(Ticker), logHandle func(string), errHandle func(error)) {
	gateway := "wss://stream.binance.com:9443/ws"
	proxyUrl := ""
	if root.UseProxy {
		proxyUrl = fmt.Sprintf("http://%s", root.ProxyUrl)
	}
	ws := websocketclient.New(gateway, proxyUrl)
	ws.OnConnectError(func(err error) {
		fmt.Printf("err: %v\n", err)
		go errHandle(err)
	})
	ws.OnDisconnected(func(err error) {
		go errHandle(err)
	})
	ws.OnSentError(func(err error) {
		go errHandle(fmt.Errorf("OnSentError: %v", err))
	})
	ws.OnConnected(func() {
		go logHandle("SubSpotTicker Connected")
		subList := []string{}
		for _, s := range symbols {
			subList = append(subList, fmt.Sprintf("%susdt@bookTicker", strings.ToLower(s)))
		}
		subData := map[string]any{
			"method": "SUBSCRIBE",
			"params": subList,
			"id":     util.GetUUID32(),
		}
		buff, _ := json.Marshal(subData)
		ws.SendTextMessage(string(buff))
		go logHandle(fmt.Sprintf("订阅币对: %v", strings.Join(symbols, "、")))
	})
	ws.OnTextMessageReceived(func(msg string) {
		m := WssSpotMsg{}
		err := json.Unmarshal([]byte(msg), &m)
		if err != nil {
			go errHandle(fmt.Errorf("msg json.Unmarshal err: %s", msg))
			return
		}
		if strings.Contains(m.Symbol, "USDT") {
			go reciveHandle(Ticker{
				Exchange: root.ExchangeName,
				Symbol:   strings.Replace(m.Symbol, "USDT", "", 1),
				Buy: Values{
					Price: util.ParseFloat(m.BuyPrice, 0),
					Size:  util.ParseFloat(m.BuySize, 0),
				},
				Sell: Values{
					Price: util.ParseFloat(m.SellPrice, 0),
					Size:  util.ParseFloat(m.SellSize, 0),
				},
				UpdateAt: root.GetTimeFloat(),
			})
		} else if m.Id != "" {
			go logHandle("订阅成功: " + m.Id)
		} else {
			go logHandle("unkown msg: " + msg)
		}

	})

	ws.OnClose(func(code int, text string) {
		// fmt.Printf("close: %v, %v\n", code, text)
		go errHandle(fmt.Errorf("close: %v, %v", code, text))
	})

	ws.Connect()

}
