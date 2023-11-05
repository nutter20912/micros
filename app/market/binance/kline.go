package binance

import (
	"encoding/json"
)

type Kline struct {
	StartTime uint64 `json:"t"` // 这K线的起始时间
	EndTime   uint64 `json:"T"` // 这K线的结束时间
	Symbol    string `json:"s"` // 交易對
	Interval  string `json:"i"` // K线间隔
	Open      string `json:"o"` // 開
	Close     string `json:"c"` // 收
	High      string `json:"h"` // 高
	Low       string `json:"l"` // 低
}

type KlieData struct {
	EventType string `json:"e"` // 事件类型
	EventTime uint64 `json:"E"` // 事件时间
	Symbol    string `json:"s"` // 交易對
	Kline     Kline  `json:"k"`
}

type KlineMessage struct {
	Stream string   `json:"stream"`
	Data   KlieData `json:"data"`
}

func newKlineMessage(message []byte) StreamMessage {
	return new(KlineMessage).parse(message)
}

func (a *KlineMessage) parse(message []byte) StreamMessage {
	json.Unmarshal(message, a)

	return a
}
