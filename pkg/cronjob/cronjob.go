package cronjob

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"todolist/internal/repository/models"

	"github.com/hibiken/asynq"
	"gopkg.in/gomail.v2"
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


func SendEmail(userID int, email, todo string) error {
	// .env dan o‘qiymiz
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPortStr := os.Getenv("SMTP_PORT")
	smtpEmail := os.Getenv("SMTP_EMAIL")
	smtpPass := os.Getenv("SMTP_PASSWORD")

	port, err := strconv.Atoi(smtpPortStr)
	if err != nil {
		return fmt.Errorf("SMTP port noto‘g‘ri: %v", err)
	}

	m := gomail.NewMessage()
	m.SetHeader("From", smtpEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Todolistdan Xabar 📌")
	m.SetBody("text/plain", fmt.Sprintf("Salom! \n\nSiz uchun yangi vazifa %s", todo))

	d := gomail.NewDialer(smtpHost, port, smtpEmail, smtpPass)

	// Email yuborish
	if err := d.DialAndSend(m); err != nil {
		log.Printf("❌ Email yuborishda xato (UserID: %d, Email: %s): %v", userID, email, err)
		return err
	}

	log.Printf("✅ Email yuborildi (UserID: %d, Email: %s)", userID, email)
	return nil
}


func HandleEmailTask(ctx context.Context, t *asynq.Task) error {
	var payload EmailPayload
	if err := json.Unmarshal(t.Payload(), &payload); err != nil {
		return err
	}
	return SendEmail(payload.UserID, payload.Email, payload.Todo)
}


func RunCronJob() {
	client := asynq.NewClient(asynq.RedisClientOpt{Addr: "localhost:6379"})
	defer client.Close()

	server := asynq.NewServer(
		asynq.RedisClientOpt{Addr: "localhost:6379"},
		asynq.Config{
			Concurrency: 10,
		},
	)

	
	go func() {
		mux := asynq.NewServeMux()
		mux.HandleFunc(TypeSendEmail, HandleEmailTask)
		if err := server.Run(mux); err != nil {
			log.Fatalf("Asynq server xato: %v", err)
		}
	}()


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

	
		for i := 0; i < len(users); i += 10 {
			end := i + 10
			if end > len(users) {
				end = len(users)
			}
			batch := users[i:end]
			for _, u := range batch {
				task, _ := NewEmailTask(int(u.ID), u.Email, "Bugun Todolist qoshishni maslahat beraman")
				info, err := client.Enqueue(task, asynq.MaxRetry(5), asynq.ProcessAt(time.Now()), asynq.Unique(15 * time.Second))
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

