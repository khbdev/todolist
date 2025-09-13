package notifactions

import (
	"context"
	"encoding/json"

	"github.com/redis/go-redis/v9"
)

type Notifier struct {
	rdb *redis.Client
}

// Notification – frontendga yuboriladigan umumiy struct
type Notification struct {
	Action string      `json:"action"` // created, updated, deleted
	Data   interface{} `json:"data"`   // todo obyekt
}

func NewNotifier(rdb *redis.Client) *Notifier {
	return &Notifier{rdb: rdb}
}

// Publish – Redis kanaliga event yuborish
func (n *Notifier) Publish(ctx context.Context, action string, data interface{}) error {
	notif := Notification{
		Action: action,
		Data:   data,
	}
	payload, err := json.Marshal(notif)
	if err != nil {
		return err
	}

	channel := "todos" // bitta umumiy kanal
	return n.rdb.Publish(ctx, channel, payload).Err()
}
