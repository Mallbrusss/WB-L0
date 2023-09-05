package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	_ "time"

	dbm "L0/DataBaseManager"
	_ "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

func main() {
	dbm.DbConnect()
	// Настройте соединение с сервером NATS Streaming.
	natsURL := "nats://localhost:4222" // Замените на ваш URL NATS Streaming сервера
	clusterID := "nat1"                // Замените на ваш Cluster ID
	clientID := "nat_client"           // Замените на ваш Client ID
	channelName := "order_chanel"      // Замените на имя вашего канала

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

	// Настройте обработчик сообщений для вашего канала.

	subscription, err := conn.Subscribe(channelName, func(msg *stan.Msg) {
		// Десериализуйте JSON-сообщение в вашу структуру данных.
		var data dbm.OrderData
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("Ошибка при десериализации JSON: %v", err)
			return
		}

		// Теперь вы можете работать с вашей структурой данных.
		fmt.Printf("Получено JSON-сообщение: %+v\n", data)

		// Запишите JSON-данные в базу данных PostgreSQL.
		if err := dbm.InsertDataToDataBase(data); err != nil {
			log.Printf("Ошибка при записи в PostgreSQL: %v", err)
		}
	})
	if err != nil {
		log.Fatalf("Ошибка при подписке на канал: %v", err)
	}
	defer func(subscription stan.Subscription) {
		err := subscription.Close()
		if err != nil {

		}
	}(subscription)

	// Ждем сигнала для завершения программы (например, Ctrl+C).
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh

	dbm.DbDisconnect()
	fmt.Println("Завершение программы...")
}
