package main

import (
	"encoding/json"
	mux2 "github.com/gorilla/mux"
	"kdf_tech_job/config"
	"kdf_tech_job/database"
	"kdf_tech_job/handlers"
	"log"
	"net/http"
)

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"status": "ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// Инициализация конфигов для проекта
	if err := config.LoadConfiguration("config.json"); err != nil {
		log.Fatal("Error while loading configurations", err)
	}

	// Инициализация соединения с базой
	if err := database.InitDatabaseConnection(); err != nil {
		log.Fatal("Error while initializing database connection", err)
	}

	r := mux2.NewRouter()

	// Обработчики
	r.HandleFunc("/healthcheck", healthCheckHandler).Methods("GET")
	r.HandleFunc("/currency/save/{date}", handlers.SaveCurrenciesHandler).Methods("GET")
	r.HandleFunc("/currency/save/{date}/{code}", handlers.SaveCurrenciesHandler).Methods("GET")

	// Слушаем наш порт
	if err := http.ListenAndServe(":"+config.Config.ListenPort, r); err != nil {
		log.Fatal("Error while serving port "+config.Config.ListenPort, err)
	}
	log.Printf("Listening on port %s\n", config.Config.ListenPort)
}
