package realtime_update

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

var newline = []byte{'\n'}

type Client struct {
	Conn *websocket.Conn
	Send chan *Payload
}

func (c *Client) WritePump(h *Hub, channelName string) {
	ticker := time.NewTicker(pingPeriod)

	defer func() {
		ticker.Stop()
		h.Unregister <- &SubscriptionType{
			ChannelName: channelName,
			Client:      c,
		}
		_ = c.Conn.Close()
		fmt.Println("conn closed")
	}()

	for {
		select {
		case message, ok := <-c.Send:

			if err := c.Conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				return
			}
			if !ok {
				_ = c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				return
			}

			jsonResponse, err := json.Marshal(message)
			if err != nil {
				return
			}

			if _, err := w.Write(jsonResponse); err != nil {
				return
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			_ = c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
