package binance

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type aggTradeData struct {
	EventType       string  `json:"e"`        // 事件类型
	EventTime       uint64  `json:"E"`        // 事件时间
	Symbol          string  `json:"s"`        // 交易對
	Price           float64 `json:"p,string"` // 成交价格
	Quantity        float64 `json:"q,string"` // 成交数量
	TransactionTime uint64  `json:"T"`        // 成交时间
	IsSell          bool    `json:"m"`        // 买方是否是做市方。如true，则此次成交是一个主动卖出单，否则是一个主动买入单。
}

type AggTradeMessage struct {
	Stream string       `json:"stream"`
	Data   aggTradeData `json:"data"`
}

func newAggTradeMessage(message []byte) StreamMessage {
	return new(AggTradeMessage).parse(message)
}

func (a *AggTradeMessage) parse(message []byte) StreamMessage {
	json.Unmarshal(message, a)

	price := fmt.Sprintf("%.1f", a.Data.Price)
	a.Data.Price, _ = strconv.ParseFloat(price, 64)

	return a
}
