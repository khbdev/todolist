package rabbitmq

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

// EmailJob – job struct
type EmailJob struct {
	Email string `json:"email"`
	Retry int    `json:"retry"`
}

// RabbitMQ struct – connection va channel
type RabbitMQ struct {
	Conn *amqp.Connection
	Ch   *amqp.Channel
}

var instance *RabbitMQ
var once sync.Once

// Singleton – faqat bitta ulanish bo‘lsin
func GetInstance() *RabbitMQ {
	once.Do(func() {
		url := os.Getenv("RABBITMQ_URL")
		conn, err := amqp.Dial(url)
		if err != nil {
			log.Fatal("RabbitMQ ulanish xatosi:", err)
		}
		ch, err := conn.Channel()
		if err != nil {
			log.Fatal("Channel yaratishda xato:", err)
		}

		// Exchange va Queue yaratish
		if err := ch.ExchangeDeclare("email_exchange", "direct", true, false, false, false, nil); err != nil {
			log.Fatal(err)
		}
		if _, err := ch.QueueDeclare("email_queue", true, false, false, false, nil); err != nil {
			log.Fatal(err)
		}
		if _, err := ch.QueueDeclare("email_dlq", true, false, false, false, nil); err != nil {
			log.Fatal(err)
		}
		if err := ch.QueueBind("email_queue", "signup", "email_exchange", false, nil); err != nil {
			log.Fatal(err)
		}

		instance = &RabbitMQ{
			Conn: conn,
			Ch:   ch,
		}
	})
	return instance
}

// Publish – jobni queue ga yuboradi
func (r *RabbitMQ) Publish(queue string, job EmailJob) error {
	body, err := json.Marshal(job)
	if err != nil {
		return err
	}

	return r.Ch.Publish(
		"email_exchange",
		queue,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// Consume – joblarni oladi va email yuboradi
func (r *RabbitMQ) Consume() {
	r.Ch.Qos(2, 0, false) // prefetch=2

	msgs, err := r.Ch.Consume("email_queue", "", false, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	for msg := range msgs {
		var job EmailJob
		if err := json.Unmarshal(msg.Body, &job); err != nil {
			log.Println("JSON parsing error:", err)
			msg.Ack(false)
			continue
		}

		err := SendEmail(job)
		if err != nil {
			job.Retry++
			if job.Retry > 3 {
				// DLQ ga yuborish
				dlqBody, _ := json.Marshal(job)
				r.Ch.Publish("", "email_dlq", false, false, amqp.Publishing{
					ContentType: "application/json",
					Body:        dlqBody,
				})
				msg.Ack(false)
				log.Println("Job DLQ ga yuborildi:", job.Email)
				continue
			}

			// Retry 5 soniyadan keyin, goroutine bilan bloklamasdan
			go func(j EmailJob) {
				time.Sleep(5 * time.Second)
				if err := r.Publish("signup", j); err != nil {
					log.Println("Retry publish error:", err)
				}
			}(job)

			msg.Ack(false)
			continue
		}

		msg.Ack(false)
		log.Println("Email yuborildi:", job.Email)
	}
}
