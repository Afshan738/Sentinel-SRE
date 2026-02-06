package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
	"github.com/streadway/amqp"
)

var ctx = context.Background()

type MonitorTask struct {
	ID  int    `json:"id"`
	URL string `json:"url"`
}

func main() {
	dbConnStr := "postgresql://AfshanQ:Afshan525@localhost:5432/sre_db?sslmode=disable"
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatalf("Failed to connect to Postgres: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatalf("Database is unreachable: %v", err)
	}

	amqpConn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer amqpConn.Close()

	ch, err := amqpConn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare("monitor_tasks", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("Failed to set QoS: %v", err)
	}

	msgs, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:16379",
		Password: "",
		DB:       0,
	})

	err = rdb.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis successfully")
	defer rdb.Close()

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var task MonitorTask
			err := json.Unmarshal(d.Body, &task)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
				d.Nack(false, false)
				continue
			}

			var statusCode int
			var statusText string
			var latencyMs int

			maxRetries := 3
			success := false

			for i := 0; i < maxRetries; i++ {
				fmt.Printf("Worker Checking URL (Attempt %d): %s\n", i+1, task.URL)

				start := time.Now()
				client := http.Client{Timeout: 5 * time.Second}
				resp, err := client.Get(task.URL)

				if err == nil && resp.StatusCode >= 200 && resp.StatusCode < 400 {
					latencyMs = int(time.Since(start).Milliseconds())
					statusCode = resp.StatusCode
					statusText = "UP"
					success = true
					resp.Body.Close()
					break
				}

				if i < maxRetries-1 {
					waitTime := time.Duration(1<<(i+1)) * time.Second
					fmt.Printf("Ping failed for %s. Retrying in %v...\n", task.URL, waitTime)
					time.Sleep(waitTime)
				}
			}

			if !success {
				statusCode = 0
				statusText = "DOWN"
				latencyMs = 0
				log.Printf("Final result: %s is DOWN", task.URL)
			} else {
				log.Printf("Final result: %s is UP | Latency: %dms", task.URL, latencyMs)
			}

			checkQuery := `INSERT INTO checks (monitor_id, status_code, latency_ms, status_text) VALUES ($1, $2, $3, $4)`
			_, err = db.Exec(checkQuery, task.ID, statusCode, latencyMs, statusText)
			if err != nil {
				log.Printf("Error saving check: %v", err)
			}

			updateQuery := `UPDATE monitors SET last_checked = NOW(), status = $1 WHERE id = $2`
			_, err = db.Exec(updateQuery, statusText, task.ID)
			if err != nil {
				log.Printf("Error updating monitors table: %v", err)
			}

			redisKey := fmt.Sprintf("monitor:%d:status", task.ID)
			err = rdb.Set(ctx, redisKey, statusText, 2*time.Hour).Err()
			if err != nil {
				log.Printf("Failed to cache status in Redis: %v", err)
			}

			d.Ack(false)
		}
	}()

	log.Printf("Guardian Worker Online")
	<-forever
}
