package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	agent "github.com/by2waysprojects/r-agent/model"
	"github.com/gorilla/websocket"
)

const (
	defaultTimeout       = 30 * time.Second
	reconnectInterval    = 5 * time.Second
	maxReconnectAttempts = 3
	testID               = "test-cluster-6"
)

type WebSocketClient struct {
	conn              *websocket.Conn
	url               string
	done              chan struct{}
	reconnectAttempts int
	handleMessageFunc func([]byte) error
}

func NewWebSocketClient(url string) (*WebSocketClient, error) {
	if url == "" {
		return nil, errors.New("URL is required")
	}

	client := &WebSocketClient{
		url:  url,
		done: make(chan struct{}),
	}

	go client.maintainConnection()

	return client, nil
}

func (c *WebSocketClient) SetMessageHandler(handler func([]byte) error) {
	c.handleMessageFunc = handler
}

func (c *WebSocketClient) maintainConnection() {
	for {
		select {
		case <-c.done:
			return
		default:
			ctx, cancel := context.WithTimeout(context.Background(), defaultTimeout)
			defer cancel()

			conn, _, err := websocket.DefaultDialer.DialContext(ctx, c.url, nil)
			if err != nil {
				log.Printf("Error connecting to WebSocket: %v", err)

				if c.reconnectAttempts < maxReconnectAttempts {
					c.reconnectAttempts++
					time.Sleep(reconnectInterval)
					continue
				}

				c.Close()
				return
			}

			c.conn = conn
			c.reconnectAttempts = 0

			msg := agent.AgentMessage{
				ID:     testID,
				Action: agent.Connect,
			}

			if err := conn.WriteJSON(msg); err != nil {
				log.Fatalf("Error sending connect message: %v", err)
			}

			c.handleMessages()
		}
	}
}

func (c *WebSocketClient) handleMessages() {
	defer c.conn.Close()

	for {
		select {
		case <-c.done:
			return
		default:
			_, message, err := c.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
					log.Printf("WebSocket error: %v", err)
				}
				return
			}

			if c.handleMessageFunc != nil {
				if err := c.handleMessageFunc(message); err != nil {
					log.Printf("Error processing message: %v", err)
				}
			}
		}
	}
}

func (c *WebSocketClient) Send(message interface{}) error {
	if c.conn == nil {
		return errors.New("not connected to WebSocket server")
	}

	msgBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("error marshaling message: %w", err)
	}

	return c.conn.WriteMessage(websocket.TextMessage, msgBytes)
}

func (c *WebSocketClient) Close() error {
	close(c.done)

	if c.conn != nil {
		err := c.conn.WriteMessage(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		)
		if err != nil {
			return err
		}
		return c.conn.Close()
	}
	return nil
}
