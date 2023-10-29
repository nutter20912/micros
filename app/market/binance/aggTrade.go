package binance

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type aggTradeData struct {
	EventType       string  `json:"e"`        // 事件类型
	EventTime       int     `json:"E"`        // 事件时间
	Price           float64 `json:"p,string"` // 成交价格
	Quantity        float64 `json:"q,string"` // 成交数量
	TransactionTime int     `json:"T"`        // 成交时间
	IsSell          bool    `json:"m"`        // 买方是否是做市方。如true，则此次成交是一个主动卖出单，否则是一个主动买入单。
}

type aggTradeMessage struct {
	Stream string       `json:"stream"`
	Data   aggTradeData `json:"data"`
}

func (a *aggTradeMessage) getResult(message []byte) []byte {
	json.Unmarshal(message, a)

	price := fmt.Sprintf("%.1f", a.Data.Price)
	a.Data.Price, _ = strconv.ParseFloat(price, 64)

	res, _ := json.Marshal(a)
	return res
}
