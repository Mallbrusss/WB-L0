package DataBaseManager

import (
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"log"
)

var dbConn *sql.DB

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

func InsertDataToDataBase(data OrderData) error {
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
	_, err = dbConn.Exec(`INSERT INTO orders (
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

func FetchDataFromDatabase(orderUID string) (OrderData, error) {
	query := `
        SELECT 
            order_uid, track_number, entry, delivery, payment, items, locale,
            internal_signature, customer_id, delivery_service, shardkey, sm_id,
            date_created, oof_shard
        FROM 
            orders
        WHERE 
            order_uid = $1
    `

	var (
		orderData    OrderData
		deliveryJSON []byte
		paymentJSON  []byte
		itemsJSON    []byte
	)

	err := dbConn.QueryRow(query, orderUID).Scan(
		&orderData.OrderUID, &orderData.TrackNumber, &orderData.Entry, &deliveryJSON, &paymentJSON,
		&itemsJSON, &orderData.Locale, &orderData.InternalSignature, &orderData.CustomerID,
		&orderData.DeliveryService, &orderData.ShardKey, &orderData.SMID, &orderData.DateCreated, &orderData.OOFShard,
	)
	if err != nil {
		log.Printf("Ошибка при выполнении SELECT-запроса: %v", err)
		return OrderData{}, err
	}

	if err := json.Unmarshal(deliveryJSON, &orderData.Delivery); err != nil {
		log.Printf("Ошибка при десериализации DeliveryData: %v", err)
		return OrderData{}, err
	}

	if err := json.Unmarshal(paymentJSON, &orderData.Payment); err != nil {
		log.Printf("Ошибка при десериализации PaymentData: %v", err)
		return OrderData{}, err
	}

	if err := json.Unmarshal(itemsJSON, &orderData.Items); err != nil {
		log.Printf("Ошибка при десериализации ItemData: %v", err)
		return OrderData{}, err
	}

	return orderData, nil
}

func Connect() {

	db, err := sql.Open("postgres", "user=manager dbname=orders password=secret host=localhost port=5432 sslmode=disable")
	if err != nil {
		log.Fatalf("Ошибка при подключении к PostgreSQL: %v", err)
	}
	dbConn = db
}

func Disconnect() {
	err := dbConn.Close()
	if err != nil {

	}
}
