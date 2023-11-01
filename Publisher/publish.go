package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/nats-io/stan.go"
)

func main() {
	natsURL := "nats://localhost:4222"
	clusterID := "nat1"
	clientID := "nat_client2"
	channelName := "order_chanel"

	conn, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Ошибка при подключении к NATS Streaming: %v", err)
	}
	defer func(conn stan.Conn) {
		err := conn.Close()
		if err != nil {
			return
		}
	}(conn)

	go func() {
		jsonData, err := os.ReadFile("model.json")
		if err != nil {
			log.Printf("Ошибка при чтении JSON-файла: %v", err)
			return
		}

		err = conn.Publish(channelName, jsonData)
		if err != nil {
			log.Printf("Ошибка при публикации сообщения: %v", err)
		} else {
			fmt.Println("Отправлена JSON-строка:", string(jsonData))
		}

		time.Sleep(2 * time.Second)

	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("Завершение программы...")
}
