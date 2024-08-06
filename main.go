package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	mux2 "github.com/gorilla/mux"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type UserDatabase struct {
	Type     string `json:"type"`
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
	Host     string `json:"host"`
	Port     string `json:"port"`
}

type ApplicationConfig struct {
	ListenPort   string       `json:"port"`
	UserDatabase UserDatabase `json:"database"`
}

type R_Currency struct {
	ID    int     `gorm:"Column:id" sql:"type:int;not null"`
	Title string  `gorm:"Column:title" sql:"type:varchar(60);not null"`
	Code  string  `gorm:"Column:code" sql:"type:varchar(3);not null"`
	Value float64 `gorm:"Column:value" sql:"type:numeric;not null"`
	Date  string  `gorm:"Column:a_date" sql:"type:date;not null"`
}

type CurrencyItem struct {
	Fullname    string  `xml:"fullname"`
	Title       string  `xml:"title"`
	Description float64 `xml:"description"`
	Quant       int     `xml:"quant"`
	Index       string  `xml:"index"`
	Change      string  `xml:"change"`
}

type CurrencyResponse struct {
	XMLName     xml.Name       `xml:"rates"`
	Generator   string         `xml:"generator"`
	Title       string         `xml:"title"`
	Link        string         `xml:"link"`
	Description string         `xml:"description"`
	Copyright   string         `xml:"copyright"`
	Date        string         `xml:"date"`
	Items       []CurrencyItem `xml:"item"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

var Config *ApplicationConfig
var db *gorm.DB

func LoadConfiguration() error {
	if _, err := os.Stat("config.json"); os.IsNotExist(err) {
		return err
	}
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	Config = &ApplicationConfig{}
	err = decoder.Decode(Config)
	if err != nil {
		return err
	}
	return nil
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := map[string]string{"status": "ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func initDatabaseConnection() (err error) {
	connectionsStr := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s",
		Config.UserDatabase.Host,
		Config.UserDatabase.Port,
		Config.UserDatabase.Username,
		Config.UserDatabase.Password,
		Config.UserDatabase.Name,
	)
	db, err = gorm.Open(postgres.Open(connectionsStr), &gorm.Config{})
	return nil
}

func saveCurrencies(currencies CurrencyResponse, date time.Time) {
	for _, currency := range currencies.Items {
		curr := R_Currency{
			Title: currency.Fullname,
			Code:  currency.Title,
			Value: currency.Description,
			Date:  date.Format("2006-01-02"),
		}
		func(currency R_Currency) {
			if err := db.Create(&currency).Error; err != nil {
				log.Printf("Error while saving currency %v: %v", currency.Title, err)
			} else {
				log.Printf("Saved currency %v", currency.Title)
			}
		}(curr)
	}
}

func getCurrencyByDateAndCode(date time.Time, code string) (R_Currency, error) {
	return R_Currency{}, nil
}

func getCurrencies(date time.Time) (CurrencyResponse, error) {
	url := "https://nationalbank.kz/rss/get_rates.cfm?fdate=" + date.Format("02-01-2006")
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error while fetching currencies from %v: %v", url, err)
		return CurrencyResponse{}, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error while reading body: %v", err)
		return CurrencyResponse{}, err
	}
	//log.Printf("Body: %v", string(body))
	var response CurrencyResponse
	err = xml.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Error while unmarshalling: %v", err)
		return CurrencyResponse{}, err
	}
	return response, nil
}

func saveCurrenciesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux2.Vars(r)
	dateStr, dateOk := vars["date"]
	codeStr, codeOk := vars["code"]

	if !dateOk {
		json.NewEncoder(w).Encode(Response{false, "date is required"})
	} else {
		// Parse the date (assuming the format is YYYY-MM-DD)
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			json.NewEncoder(w).Encode(Response{false, "invalid date format"})
			return
		}

		if !codeOk {
			currencyResp, err := getCurrencies(date)
			if err != nil {
				log.Printf("Error getting currencies: %v", err)
				json.NewEncoder(w).Encode(Response{false, "can't retrieve currencies"})
				return
			}
			saveCurrencies(currencyResp, date)
			json.NewEncoder(w).Encode(Response{true, "saving currency"})
		} else {
			currency, err := getCurrencyByDateAndCode(date, codeStr)
			if err != nil {
				json.NewEncoder(w).Encode(Response{false, "can't retrieve currency by code and date"})
			}
			json.NewEncoder(w).Encode(Response{true, "getting currencies" + currency.Title})
		}
	}

	// Set the Content-Type header to application/json
	w.Header().Set("Content-Type", "application/json")

	// Write the response as JSON
}

func main() {
	if err := LoadConfiguration(); err != nil {
		panic(err)
	}

	err := initDatabaseConnection()
	if err := initDatabaseConnection(); err != nil {
		panic(err)
	}

	r := mux2.NewRouter()

	r.HandleFunc("/healthCheck", healthCheckHandler).Methods("GET")
	r.HandleFunc("/currency/save/{date}", saveCurrenciesHandler).Methods("GET")
	r.HandleFunc("/currency/save/{date}/{code}", func(w http.ResponseWriter, r *http.Request) {})

	err = http.ListenAndServe(":8080", r)
	log.Println("Listening on port " + Config.ListenPort)
	if err != nil {
		panic(err)
	}
}
