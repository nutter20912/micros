package binance

import "encoding/json"

type MiniTickerData struct {
	EventType string `json:"e"` // 事件类型
	EventTime int    `json:"E"` // 事件时间
	Symbol    string `json:"s"` // 交易對

	Close string `json:"c"` // 最新成交价格
	Open  string `json:"o"` // 24小时前开始第一笔成交价格
	High  string `json:"h"` // 24小时内最高成交价
	Low   string `json:"l"` // 24小时内最低成交价
}

type MiniTickerArrMessage struct {
	Stream string           `json:"stream"`
	Data   []MiniTickerData `json:"data"`
}

func newMiniTickerArrMessage(message []byte) StreamMessage {
	return new(MiniTickerArrMessage).parse(message)
}

func (a *MiniTickerArrMessage) parse(message []byte) StreamMessage {
	json.Unmarshal(message, a)

	var newData []MiniTickerData

	for _, v := range a.Data {
		if v.Symbol == "BTCUSDT" {
			newData = append(newData, v)
		}
	}

	a.Data = newData

	return a
}
