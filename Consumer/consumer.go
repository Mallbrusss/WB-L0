package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	_ "time"

	dbm "L0/DataBaseManager"
	_ "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

var natConn stan.Conn
var nutSubscription stan.Subscription

func Consume() {
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
	natConn = conn

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
	nutSubscription = subscription
	if err != nil {
		log.Fatalf("Ошибка при подписке на канал: %v", err)
	}
}

func Disconnect() {
	handleUnsubscribe()
	handleDisconnect()
}

func handleUnsubscribe() {
	err := nutSubscription.Close()
	if err != nil {

	}
}

func handleDisconnect() {
	err := natConn.Close()
	if err != nil {

	}
}
