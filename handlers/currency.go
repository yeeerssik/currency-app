package handlers

import (
	"encoding/json"
	mux2 "github.com/gorilla/mux"
	"kdf_tech_job/api"
	"kdf_tech_job/database"
	"log"
	"net/http"
	"time"
)

type BaseResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

type CurrencyData struct {
	Title *string    `json:"title"`
	Code  *string    `json:"code"`
	Value *float64   `json:"value"`
	Date  *time.Time `json:"date"`
}

type GetCurrencyResponse struct {
	BaseResponse
	Currency []*CurrencyData `json:"data"`
}

func SaveCurrenciesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux2.Vars(r)
	dateStr, dateOk := vars["date"]
	codeStr, codeOk := vars["code"]

	if !dateOk {
		json.NewEncoder(w).Encode(BaseResponse{false, "date is required"})
	} else {
		// Парсим дату
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			json.NewEncoder(w).Encode(BaseResponse{false, "invalid date format"})
			return
		}

		if !codeOk {
			currencyResp, err := api.GetApiCurrencies(date)
			if err != nil {
				log.Printf("Error getting currencies: %v", err)
				json.NewEncoder(w).Encode(BaseResponse{false, "can't retrieve currencies"})
				return
			}
			var currencies []database.Currency
			for _, currencyItem := range currencyResp.Items {
				currencies = append(currencies, database.Currency{
					Title: &currencyItem.Fullname,
					Code:  &currencyItem.Title,
					Value: &currencyItem.Description,
					Date:  &date,
				})
			}
			database.SaveCurrencies(currencies)

			json.NewEncoder(w).Encode(BaseResponse{true, "saved currencies"})
		} else {
			currency, err := database.GetCurrenciesByDateAndCode(date, codeStr)
			var currencies []*CurrencyData
			if err != nil {
				json.NewEncoder(w).Encode(BaseResponse{false, "can't retrieve currency by code and date"})
			}
			for _, currencyData := range currency {
				currencies = append(currencies, &CurrencyData{
					Title: currencyData.Title,
					Code:  currencyData.Code,
					Value: currencyData.Value,
					Date:  currencyData.Date,
				})
			}
			json.NewEncoder(w).Encode(GetCurrencyResponse{
				BaseResponse{true, "fetched currencies"},
				currencies,
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
}
