package cronjob

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"
	"todolist/internal/repository/models"

	"github.com/hibiken/asynq"
	"gorm.io/gorm"

	"net/smtp"
)

type EmailTaskPayload struct {
    Email string
    Msg   string
}

func NewEmailTask(email, msg string) (*asynq.Task, error) {
    payload, err := json.Marshal(EmailTaskPayload{
        Email: email,
        Msg:   msg,
    })
    if err != nil {
        return nil, err
    }
    return asynq.NewTask("email:send", payload), nil
}

// Email yuborish funksiyasi (Gmail SMTP)
func sendEmailGmail(to, msg string) error {
    from := os.Getenv("SMTP_EMAIL")
    password := os.Getenv("SMTP_PASSWORD")
    smtpHost := os.Getenv("SMTP_HOST")
    smtpPort := os.Getenv("SMTP_PORT") // e.g., "587"

    body := "Subject: TodoList Update\n\n" + msg
    auth := smtp.PlainAuth("", from, password, smtpHost)
    return smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(body))
}

// processEmails -> db qabul qiladi
func processEmails(db *gorm.DB) {
    ctx := context.Background()
    users, err := models.GetAllUsers(db) // db uzatyapmiz
    if err != nil {
        log.Println("DB xato:", err)
        return
    }

    for _, u := range users {
        redisKey := fmt.Sprintf("email_sent:%s", u.Email)

        exists, err := RedisClient.Exists(ctx, redisKey).Result()
        if err != nil {
            log.Println("Redis xato:", err)
            continue
        }
        if exists > 0 {
            continue
        }

        task, err := NewEmailTask(u.Email, "Welcome to TodoList!")
        if err != nil {
            log.Println("Task yaratishda xato:", err)
            continue
        }

        _, err = AsynqClient.Enqueue(task)
        if err != nil {
            log.Println("Task queue-ga yuborishda xato:", err)
            continue
        }

        // Redis TTL = 1 soat
        err = RedisClient.Set(ctx, redisKey, "sent", time.Hour).Err()
        if err != nil {
            log.Println("Redis set xato:", err)
        }

        log.Println("✅ Email queued for:", u.Email)
    }
}

// RunHourlyCronJob -> db qabul qiladi
func RunHourlyCronJob(db *gorm.DB) {
    go func() {
        ticker := time.NewTicker(5 * time.Second)
        for range ticker.C {
            processEmails(db)
        }
    }()
}


// Asynq Server – emaillarni orqa fonda ishlash
func StartWorker() {
    srv := asynq.NewServer(
        asynq.RedisClientOpt{Addr: "localhost:6379"},
        asynq.Config{
            Concurrency: 10,
        },
    )

    mux := asynq.NewServeMux()
    mux.HandleFunc("email:send", func(ctx context.Context, t *asynq.Task) error {
        var p EmailTaskPayload
        if err := json.Unmarshal(t.Payload(), &p); err != nil {
            return err
        }
        log.Printf("Sending email to %s", p.Email)
        return sendEmailGmail(p.Email, p.Msg)
    })

    if err := srv.Run(mux); err != nil {
        log.Fatal(err)
    }
}
