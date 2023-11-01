package HttpServer

import (
	cache "L0/Cache"
	"L0/DataBaseManager"
	"encoding/json"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"time"
)


func Serv() {
	myCache := cache.NewCache()

	http.HandleFunc("/getOrderData", func(w http.ResponseWriter, r *http.Request) {
		orderUID := r.URL.Query().Get("orderUID")

		cachedData, exists := myCache.Get(orderUID)
		if exists {
			jsonData, err := json.Marshal(cachedData)
			if err != nil {
				http.Error(w, "Error converting data to JSON", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
			return
		}

		orderData, err := DataBaseManager.FetchDataFromDatabase(orderUID)
		if err != nil {
			http.Error(w, "Error fetching data from the database", http.StatusInternalServerError)
			return
		}

		myCache.Set(orderUID, orderData, 5*time.Minute)

		jsonData, err := json.Marshal(orderData)
		if err != nil {
			http.Error(w, "Error converting data to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	fs := http.FileServer(http.Dir("./Static"))
	http.Handle("/", fs)

	log.Println("Сервер запущен на localhost:8080")
	_ = http.ListenAndServe(":8080", nil)
}
