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
	natsURL := "nats://localhost:4222" // Замените на ваш URL NATS Streaming сервера
	clusterID := "nat1"                // Замените на ваш Cluster ID
	clientID := "nat_client"           // Замените на ваш Client ID
	channelName := "order_chanel"      // Замените на имя вашего канала

	conn, err := stan.Connect(clusterID, clientID, stan.NatsURL(natsURL))
	if err != nil {
		log.Fatalf("Ошибка при подключении к NATS Streaming: %v", err)
	}
	natConn = conn

	subscription, err := conn.Subscribe(channelName, func(msg *stan.Msg) {
		var data dbm.OrderData
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("Ошибка при десериализации JSON: %v", err)
			return
		}

		fmt.Printf("Получено JSON-сообщение: %+v\n", data)

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
		return
	}
}

func handleDisconnect() {
	err := natConn.Close()
	if err != nil {
		return
	}
}
