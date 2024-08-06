package database

import (
	"log"
	"time"
)

type Currency struct {
	ID    *int       `gorm:"Column:id" sql:"type:int;not null"`
	Title *string    `gorm:"Column:title" sql:"type:varchar(60);not null"`
	Code  *string    `gorm:"Column:code" sql:"type:varchar(3);not null"`
	Value *float64   `gorm:"Column:value" sql:"type:numeric;not null"`
	Date  *time.Time `gorm:"Column:a_date" sql:"type:date;not null"`
}

func (Currency) TableName() string {
	return "r_currency"
}

func SaveCurrencies(currencies []Currency) {
	for _, currency := range currencies {
		go func(currency Currency) {
			if err := db.Create(&currency).Error; err != nil {
				log.Printf("Error while saving currency %v: %v", currency.Title, err)
			} else {
				log.Printf("Saved currency %v", *(currency.Title))
			}
		}(currency)
	}
}

func GetCurrenciesByDateAndCode(date time.Time, code string) (currencies []*Currency, err error) {
	if find := db.Where("a_date = ? and code = ?", date, code).Find(&currencies).Error; find != nil {
		return nil, find
	}
	return currencies, nil
}
