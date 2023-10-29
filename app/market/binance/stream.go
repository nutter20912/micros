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
		m := &depthMessage{}
		message = m.getResult(message)

	case bytes.Contains(message, []byte("aggTrade")):
		m := &aggTradeMessage{}
		message = m.getResult(message)

	case bytes.Contains(message, []byte("!miniTicker@arr")):
		m := &MiniTickerArrMessage{}
		message = m.getResult(message)
	}

	return message, nil
}

func (s *Stream) Close() {
	s.conn.Close()
}
