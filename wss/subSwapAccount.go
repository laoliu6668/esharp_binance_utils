package binance_wss

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	htx "github.com/laoliu6668/esharp_binance_utils"
	"github.com/laoliu6668/esharp_binance_utils/apis"
	"github.com/laoliu6668/esharp_binance_utils/util"
	"github.com/laoliu6668/esharp_binance_utils/util/websocketclient"
)

// 订阅期货账户变化
// 余额+持仓+订单
func SubSwapAccount(reciveHandle func(ReciveBalanceMsg), logHandle func(string), errHandle func(error)) {

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
	if htx.UseProxy {
		go logHandle(fmt.Sprintf("proxyUrl: %v\n", htx.ProxyUrl))
		proxyUrl = fmt.Sprintf("http://%s", htx.ProxyUrl)
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
		sub := map[string]any{
			"method": "REQUEST",
			"params": []string{
				fmt.Sprintf("%s@account", listenKey),
				fmt.Sprintf("%s@balance", listenKey),
			},
			"id": 12,
		}
		buf, _ := json.Marshal(sub)
		ws.SendTextMessage(string(buf))
		logHandle(fmt.Sprintf("send sub msg: %s", string(buf)))
	})
	ws.OnTextMessageReceived(func(message string) {
		fmt.Printf("message: %v\n", message)
	})
	ws.OnBinaryMessageReceived(func(message []byte) {
		fmt.Printf("message: %v\n", string(message))
		type Msg struct {
			Action string         `json:"action"`
			Ch     string         `json:"ch"`
			Code   int            `json:"code"`
			Data   map[string]any `json:"data"`
		}
		msg := Msg{}
		err := json.Unmarshal([]byte(message), &msg)
		if err != nil {
			go errHandle(fmt.Errorf("decode: %v", err))
			return
		}
		if msg.Action == "ping" {
			type pingTs struct {
				Ts int64 `json:"ts"`
			}
			type pingRes struct {
				Action string `json:"action"`
				Data   pingTs `json:"data"`
			}
			pingRet := &pingRes{}
			json.Unmarshal([]byte(message), pingRet)
			pong := fmt.Sprintf(`{"action":"pong","data":{"ts":%d}}`, pingRet.Data.Ts)
			// 收到ping 回复pong
			ws.SendTextMessage(pong)
		} else if msg.Action == "push" && strings.Contains(msg.Ch, "accounts.update") {

			type Data struct {
				Currency    string      `json:"currency"`
				AccountId   int64       `json:"accountId"`
				Balance     json.Number `json:"balance"`
				Available   json.Number `json:"available"`
				AccountType string      `json:"accountType"`
				// SeqNum      int64       `json:"seqNum"`
			}
			type TickerRes struct {
				Data Data `json:"data"`
			}
			res := TickerRes{}
			json.Unmarshal([]byte(message), &res)
			if res.Data.AccountType == "trade" {
				a, _ := res.Data.Available.Float64()
				b, _ := res.Data.Balance.Float64()
				go reciveHandle(ReciveBalanceMsg{
					Exchange:  htx.ExchangeName,
					Symbol:    strings.ToUpper(res.Data.Currency),
					AccountId: res.Data.AccountId,
					Available: a,
					Balance:   b,
				})
			}
		} else if msg.Action == "req" {
			if msg.Ch == "auth" && msg.Code == 200 {
				// 订阅账户信息
				subAccountUpdateMp := map[string]any{
					"action": "sub",
					"cid":    util.GetUUID32(),
					"ch":     "accounts.update#0",
				}
				bf, _ := json.Marshal(subAccountUpdateMp)
				// fmt.Printf("sub: %v\n", string(bf))
				go logHandle(fmt.Sprintf("sub: %v\n", string(bf)))
				ws.SendTextMessage(string(bf))
			}
		} else if msg.Action == "sub" {
			go logHandle("## SubAccountUpdate sub success\n")
			// if msg.Code != 200 {
			// }
		} else {
			// fmt.Printf("unknown message: %v\n", string(message))
			go logHandle(fmt.Sprintf("unknown message: %v\n", string(message)))
		}
	})

	ws.OnClose(func(code int, text string) {
		// fmt.Printf("close: %v, %v\n", code, text)
		go errHandle(fmt.Errorf("close: %v, %v", code, text))
	})

	ws.Connect()

}
