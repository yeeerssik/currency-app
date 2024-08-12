//   Currency Api:
//    version: 0.1
//    title: Currency Api
//   Schemes: http, https
//   Host:
//   BasePath: /currency
//      Consumes:
//      - application/json
//   Produces:
//   - application/json
//   swagger:meta
package main

import (
	"encoding/json"
	"kdf_tech_job/config"
	"kdf_tech_job/database"
	"kdf_tech_job/handlers"
	"log"
	"net/http"

	"github.com/go-openapi/runtime/middleware"
	mux2 "github.com/gorilla/mux"
)

// swagger:operation GET /healthcheck
// Health Check
//
// ---
// responses:
//
//	200: CommonSuccess
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
	r.HandleFunc("/currency/{date}/{code}", handlers.GetCurrenciesHandler).Methods("GET")
	r.HandleFunc("/currency/{date}", handlers.GetCurrenciesHandler).Methods("GET")

	// swagger
	r.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))
	opts := middleware.SwaggerUIOpts{SpecURL: "swagger.yaml"}
	sh := middleware.SwaggerUI(opts, nil)
	r.Handle("/docs", sh)

	// Документация
	opts1 := middleware.RedocOpts{SpecURL: "swagger.yaml", Path: "doc"}
	sh1 := middleware.Redoc(opts1, nil)
	r.Handle("/doc", sh1)

	// Слушаем наш порт
	if err := http.ListenAndServe(":"+config.Config.ListenPort, r); err != nil {
		log.Fatal("Error while serving port "+config.Config.ListenPort, err)
	}
	log.Printf("Listening on port %s\n", config.Config.ListenPort)
}
