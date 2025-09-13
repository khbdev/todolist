package notifactions

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/redis/go-redis/v9"
)

type Subscriber struct {
	rdb *redis.Client
}

func NewSubscriber(rdb *redis.Client) *Subscriber {
	return &Subscriber{rdb: rdb}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // ⚡️ barcha originlarga ruxsat
	},
}

// HandleWS – websocket ulanishni boshqaradi
func (s *Subscriber) HandleWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Upgrade error:", err)
		return
	}
	defer conn.Close()

	ctx := context.Background()
	pubsub := s.rdb.Subscribe(ctx, "todos") // umumiy kanal
	defer pubsub.Close()

	ch := pubsub.Channel()

	for msg := range ch {
		// Redis’dan kelgan xabarni to‘g‘ridan-to‘g‘ri websocketga yuboramiz
		err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload))
		if err != nil {
			fmt.Println("WebSocket write error:", err)
			return
		}
	}
}
