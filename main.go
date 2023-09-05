package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	_ "time"

	_ "github.com/lib/pq"
	_ "github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
)

type DeliveryData struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}
type PaymentData struct {
	Transaction  string `json:"transaction"`
	RequestID    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int64  `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}
type ItemData struct {
	ChrtID      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	RID         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmID        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}
type OrderData struct {
	OrderUID          string       `json:"order_uid"`
	TrackNumber       string       `json:"track_number"`
	Entry             string       `json:"entry"`
	Delivery          DeliveryData `json:"delivery"`
	Payment           PaymentData  `json:"payment"`
	Items             []ItemData   `json:"items"`
	Locale            string       `json:"locale"`
	InternalSignature string       `json:"internal_signature"`
	CustomerID        string       `json:"customer_id"`
	DeliveryService   string       `json:"delivery_service"`
	ShardKey          string       `json:"shardkey"`
	SMID              int          `json:"sm_id"`
	DateCreated       string       `json:"date_created"`
	OOFShard          string       `json:"oof_shard"`
}

func main() {
	// Настройте соединение с сервером NATS Streaming.
	natsURL := "nats://localhost:4222" // Замените на ваш URL NATS Streaming сервера
	clusterID := "nat1"                // Замените на ваш Cluster ID
	clientID := "nat_client"           // Замените на ваш Client ID

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
	channelName := "order_chanel" // Замените на имя вашего канала

	db, err := sql.Open("postgres", "user=manager dbname=orders password=secret host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatalf("Ошибка при подключении к PostgreSQL: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {

		}
	}(db)

	subscription, err := conn.Subscribe(channelName, func(msg *stan.Msg) {
		// Десериализуйте JSON-сообщение в вашу структуру данных.
		var data OrderData
		if err := json.Unmarshal(msg.Data, &data); err != nil {
			log.Printf("Ошибка при десериализации JSON: %v", err)
			return
		}

		// Теперь вы можете работать с вашей структурой данных.
		fmt.Printf("Получено JSON-сообщение: %+v\n", data)

		// Запишите JSON-данные в базу данных PostgreSQL.
		if err := insertDataToDataBase(db, data); err != nil {
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

	fmt.Println("Завершение программы...")
}

func insertDataToDataBase(db *sql.DB, data OrderData) error {
	deliveryJSON, err := json.Marshal(data.Delivery)
	if err != nil {
		return err
	}
	paymentJSON, err := json.Marshal(data.Payment)
	if err != nil {
		return err
	}

	itemsJSON, err := json.Marshal(data.Items)
	if err != nil {
		return err
	}
	_, err = db.Exec(`
        INSERT INTO orders (
            order_uid, track_number, entry, delivery, payment, items, locale,
            internal_signature, customer_id, delivery_service, shardkey, sm_id,
            date_created, oof_shard
        ) VALUES (
            $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14
        )`,
		data.OrderUID, data.TrackNumber, data.Entry, deliveryJSON, paymentJSON,
		itemsJSON, data.Locale, data.InternalSignature, data.CustomerID,
		data.DeliveryService, data.ShardKey, data.SMID, data.DateCreated,
		data.OOFShard)
	if err != nil {
		return err
	}
	return nil
}
