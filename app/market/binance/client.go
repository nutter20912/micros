package binance

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	SPOT_URL = "wss://stream.binance.com:9443/stream"
)

type Client struct {
	Ctx context.Context
}

type RequestMessage struct {
	Method string   `json:"method"`
	Params []string `json:"params"`
	Id     int64    `json:"id"`
}

func NewClient(ctx context.Context) *Client {
	return &Client{Ctx: ctx}
}

func (c *Client) MiniTickersStream() (*Stream, error) {
	conn, response, err := websocket.DefaultDialer.Dial(SPOT_URL, http.Header{})
	if err != nil {
		fmt.Println("連接失敗:", err)
		return nil, err
	}

	if response.StatusCode != http.StatusSwitchingProtocols {
		fmt.Println("握手失敗:", response.Status)
		return nil, err
	}

	mstr, _ := json.Marshal(RequestMessage{
		Method: "SUBSCRIBE",
		Params: []string{"!miniTicker@arr"},
		Id:     time.Now().Unix(),
	})

	conn.WriteMessage(websocket.TextMessage, mstr)

	return &Stream{conn: conn}, nil
}

func (c *Client) Stream(symbol string) (*Stream, error) {
	conn, response, err := websocket.DefaultDialer.Dial(SPOT_URL, http.Header{})
	if err != nil {
		fmt.Println("連接失敗:", err)
		return nil, err
	}

	if response.StatusCode != http.StatusSwitchingProtocols {
		fmt.Println("握手失敗:", response.Status)
		return nil, err
	}

	symbol = strings.ToLower(symbol)

	mstr, _ := json.Marshal(RequestMessage{
		Method: "SUBSCRIBE",
		Params: []string{
			fmt.Sprintf("%s@kline_1m", symbol),
			fmt.Sprintf("%s@depth10@100ms", symbol),
			fmt.Sprintf("%s@aggTrade", symbol),
		},
		Id: time.Now().Unix(),
	})

	conn.WriteMessage(websocket.TextMessage, mstr)

	return &Stream{conn: conn}, nil
}
