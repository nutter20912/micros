package binance

import (
	"encoding/json"
	"fmt"
	"micros/helper"
	"sort"
	"strconv"
)

type depthData struct {
	Bids [][]string `json:"bids"` //委買
	Asks [][]string `json:"asks"` //委賣
}

type DepthMessage struct {
	Stream string    `json:"stream"`
	Data   depthData `json:"data"`
}

func newDepthMessage(message []byte) StreamMessage {
	return new(DepthMessage).parse(message)
}

func (d *DepthMessage) getMergeDepth(data [][]string) [][]string {
	merge := map[string]string{}

	for _, value := range data {
		priceF64, _ := strconv.ParseFloat(value[0], 64)
		priceStr := fmt.Sprintf("%.1f", priceF64)

		if total, ok := merge[priceStr]; ok {
			merge[priceStr] = helper.Add(total, value[1])
		} else {
			merge[priceStr] = value[1]
		}
	}

	result := [][]string{}

	for price, total := range merge {
		result = append(result, []string{price, total})
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i][0] < result[j][0]
	})

	return result
}

// 合併結果至整數
func (d *DepthMessage) parse(message []byte) StreamMessage {
	json.Unmarshal(message, d)

	d.Data.Bids = d.getMergeDepth(d.Data.Bids)
	d.Data.Asks = d.getMergeDepth(d.Data.Asks)

	return d
}
