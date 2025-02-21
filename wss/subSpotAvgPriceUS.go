package wss

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"

	root "github.com/laoliu6668/esharp_binance_utils"
	"github.com/laoliu6668/esharp_binance_utils/util/websocketclient"
)

type WssSpotAvgPriceMsgUS struct {
	Id                   string `json:"id"`
	EventType            string `json:"e"` // Event type
	EventTime            int64  `json:"E"` // Event time
	Symbol               string `json:"s"` // 交易对
	Coin                 string `json:"coin"`
	AveragePriceInterval string `json:"i"` // Average price interval
	AveragePrice         string `json:"w"` // Average price
	LastTradeTime        int64  `json:"T"` // Last trade time
}

func SubSpotAvgPriceUS(symbols []string, reciveHandle func(WssSpotAvgPriceMsgUS), logHandle func(string), errHandle func(error)) {
	gateway := "wss://stream.binance.us:9443/ws"
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
			subList = append(subList, fmt.Sprintf("%susdt@ticker", strings.ToLower(s)))
		}
		subData := map[string]any{
			"method": "SUBSCRIBE",
			"params": subList,
			"id":     rand.Intn(99999999999),
		}
		buff, _ := json.Marshal(subData)
		ws.SendTextMessage(string(buff))
		go logHandle(fmt.Sprintf("订阅币对: %v", strings.Join(symbols, "、")))
	})
	ws.OnTextMessageReceived(func(msg string) {
		m := WssSpotAvgPriceMsgUS{}
		err := json.Unmarshal([]byte(msg), &m)
		if err != nil {
			go errHandle(fmt.Errorf("msg json.Unmarshal err: %s", msg))
			return
		}
		m.Coin = strings.Replace(m.Symbol, "USDT", "", 1)
		if strings.Contains(m.Symbol, "USDT") {
			go reciveHandle(m)
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
