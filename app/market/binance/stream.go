package binance

import (
	"bytes"
	"fmt"

	"github.com/gorilla/websocket"
)

type Stream struct {
	conn *websocket.Conn
}

func (s *Stream) ReadMessage() ([]byte, error) {
	_, message, err := s.conn.ReadMessage()
	if err != nil {
		fmt.Println("讀取消息失敗:", err)
		return nil, err
	}

	switch {
	case bytes.Contains(message, []byte("depth")):
		message = new(depthMessage).getResult(message)

	case bytes.Contains(message, []byte("aggTrade")):
		message = new(AggTradeMessage).getResult(message)

	case bytes.Contains(message, []byte("!miniTicker@arr")):
		message = new(MiniTickerArrMessage).getResult(message)
	}

	return message, nil
}

func (s *Stream) Close() {
	s.conn.Close()
}
