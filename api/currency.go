package api

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

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

func GetApiCurrencies(date time.Time) (CurrencyResponse, error) {
	url := fmt.Sprintf("https://nationalbank.kz/rss/get_rates.cfm?fdate=%s",
		date.Format("02-01-2006"),
	)
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
	var response CurrencyResponse
	err = xml.Unmarshal(body, &response)
	if err != nil {
		log.Printf("Error while unmarshalling: %v", err)
		return CurrencyResponse{}, err
	}
	return response, nil
}
