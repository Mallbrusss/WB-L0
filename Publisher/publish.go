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
	// Настройте соединение с сервером NATS Streaming.
	natsURL := "nats://localhost:4222"
	clusterID := "nat1"
	clientID := "nat_client2"
	channelName := "order_chanel"
	// Создайте подключение к серверу NATS Streaming.
	conn, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Ошибка при подключении к NATS Streaming: %v", err)
	}
	defer func(conn stan.Conn) {
		err := conn.Close()
		if err != nil {

		}
	}(conn)

	go func() {
		// Считайте JSON-данные из файла.
		jsonData, err := os.ReadFile("model.json")
		if err != nil {
			log.Printf("Ошибка при чтении JSON-файла: %v", err)
			return
		}

		// Отправьте JSON-строку в канал.
		err = conn.Publish(channelName, jsonData)
		if err != nil {
			log.Printf("Ошибка при публикации сообщения: %v", err)
		} else {
			fmt.Println("Отправлена JSON-строка:", string(jsonData))
		}

		time.Sleep(2 * time.Second) // Ожидание перед отправкой следующего сообщения

	}()

	// Ждем сигнала для завершения программы (например, Ctrl+C).
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	fmt.Println("Завершение программы...")
}
