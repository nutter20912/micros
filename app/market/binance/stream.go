package binance

import (
	"bytes"
	"fmt"

	"github.com/gorilla/websocket"
)

type Stream struct {
	conn *websocket.Conn
}

func (s *Stream) ReadMessage() (StreamMessage, error) {
	_, message, err := s.conn.ReadMessage()
	if err != nil {
		fmt.Println("讀取消息失敗:", err)
		return nil, err
	}

	return MessageFactory(message), nil
}

func (s *Stream) Close() {
	s.conn.Close()
}

type StreamMessage interface {
	parse([]byte) StreamMessage
}

func MessageFactory(message []byte) StreamMessage {
	switch {
	case bytes.Contains(message, []byte("depth")):
		//return newDepthMessage(message)

	case bytes.Contains(message, []byte("aggTrade")):
		return newAggTradeMessage(message)

	case bytes.Contains(message, []byte("kline")):
		return newKlineMessage(message)

	case bytes.Contains(message, []byte("!miniTicker@arr")):
		return newMiniTickerArrMessage(message)
	}

	return nil
}
