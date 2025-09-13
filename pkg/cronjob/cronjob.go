package cronjob

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hibiken/asynq"
	"todolist/internal/repository/models"
)

const TypeSendEmail = "email:send"

type EmailPayload struct {
	UserID int
	Email  string
	Todo   string
}

func NewEmailTask(userID int, email, todo string) (*asynq.Task, error) {
	payload := EmailPayload{UserID: userID, Email: email, Todo: todo}
	data, _ := json.Marshal(payload)
	return asynq.NewTask(TypeSendEmail, data), nil
}

// SendEmail hozircha faqat print
func SendEmail(userID int, email, todo string) error {
	fmt.Printf("✅ [Test] UserID: %d ga email yuborilyapti: %s, Todo: %s\n", userID, email, todo)
	return nil
}

// Worker
func HandleEmailTask(ctx context.Context, t *asynq.Task) error {
	var payload EmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	return SendEmail(payload.UserID, payload.Email, payload.Todo)
}

// Cron ishga tushadi
func RunCronJob() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379"})
	defer client.Close()

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{
			Concurrency: 10,
		},
	)

	// Worker gorutinda
	go func() {
		mux := asynq.NewServeMux()
		mux.HandleFunc(TypeSendEmail, HandleEmailTask)
		if err := server.Run(mux); err != nil {
			log.Fatalf("Asynq server xato: %v", err)
		}
	}()

	// Har Interval ishlaydi
	ticker := time.NewTicker(Interval)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("⏱ Cron job ishga tushdi")

		// DB dan userlarni olish
		var users []models.User
		if err := DB.Find(&users).Error; err != nil {
			log.Println("DB dan user olish xato:", err)
			continue
		}

		// Batch 10 user
		for i := 0; i < len(users); i += 10 {
			end := i + 10
			if end > len(users) {
				end = len(users)
			}
			batch := users[i:end]
			for _, u := range batch {
				task, _ := NewEmailTask(int(u.ID), u.Email, "Bugungi todo ro'yxati")
				info, err := client.Enqueue(task, asynq.MaxRetry(5), asynq.ProcessAt(time.Now()))
				if err != nil {
					log.Println("Task enqueue xato:", err)
				} else {
					log.Println("Task enqueued ID:", info.ID)
				}
			}
			time.Sleep(50 * time.Millisecond)
		}
	}
}
