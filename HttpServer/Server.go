package HttpServer

import (
	"L0/DataBaseManager" // Замените на путь к вашему пакету
	"encoding/json"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

// Инициализируйте базу данных здесь

func Serv() {
	// Настройте маршрут и обработчик HTTP-запросов.
	http.HandleFunc("/getOrderData", func(w http.ResponseWriter, r *http.Request) {
		// Получите значение параметра orderUID из запроса.
		orderUID := r.URL.Query().Get("orderUID")

		// Вызовите функцию FetchDataFromDatabase для получения данных из базы данных.
		orderData, err := DataBaseManager.FetchDataFromDatabase(orderUID)
		if err != nil {
			http.Error(w, "Ошибка при получении данных из базы данных", http.StatusInternalServerError)
			return
		}

		// Преобразуйте данные в JSON и отправьте их в ответе.
		jsonData, err := json.Marshal(orderData)
		if err != nil {
			http.Error(w, "Ошибка при преобразовании данных в JSON", http.StatusInternalServerError)
			return
		}

		// Установите заголовок Content-Type и отправьте JSON-данные.
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonData)
	})

	// Запустите HTTP-сервер на порту 8080.,
	fs := http.FileServer(http.Dir("./Static"))
	http.Handle("/", fs)

	log.Println("Сервер запущен на localhost:8080")
	_ = http.ListenAndServe(":8080", nil)
}
