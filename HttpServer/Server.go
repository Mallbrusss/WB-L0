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

// Инициализируйте базу данных здесь

func Serv() {
	// Create an instance of the cache.
	myCache := cache.NewCache()

	// Set up the HTTP route and request handler.
	http.HandleFunc("/getOrderData", func(w http.ResponseWriter, r *http.Request) {
		// Get the value of the orderUID parameter from the request.
		orderUID := r.URL.Query().Get("orderUID")

		// Try to get data from the cache.
		cachedData, exists := myCache.Get(orderUID)
		if exists {
			// If data is in the cache, send it in the response.
			jsonData, err := json.Marshal(cachedData)
			if err != nil {
				http.Error(w, "Error converting data to JSON", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonData)
			return
		}

		// If data is not in the cache, call the FetchDataFromDatabase function.
		orderData, err := DataBaseManager.FetchDataFromDatabase(orderUID)
		if err != nil {
			http.Error(w, "Error fetching data from the database", http.StatusInternalServerError)
			return
		}

		// Save the retrieved data in the cache for a specified duration.
		myCache.Set(orderUID, orderData, 5*time.Minute) // Example: cache will be stored for 1 minute.

		// Convert the data to JSON and send it in the response.
		jsonData, err := json.Marshal(orderData)
		if err != nil {
			http.Error(w, "Error converting data to JSON", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	// Запустите HTTP-сервер на порту 8080.
	fs := http.FileServer(http.Dir("./Static"))
	http.Handle("/", fs)

	log.Println("Сервер запущен на localhost:8080")
	_ = http.ListenAndServe(":8080", nil)
}
