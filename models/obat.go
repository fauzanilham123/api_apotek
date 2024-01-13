package models

import "time"

type (
	Drug struct {
		ID        		uint     		`gorm:"primary_key" json:"id"`
		Kode        	string      	`gorm:"unique" json:"kode"`
		Name        	string      	`json:"obat"`
		ExpiredDate     string      	`json:"expired_date"`
		Jumlah        	int      		`json:"jumlah"`
		HargaPerUnit    int      		`json:"harga_per_unit"`
		Flag			int 			`json:"flag"`
		CreatedAt 		time.Time 		`json:"created_at"`
		UpdatedAt 		time.Time 		`json:"updated_at"`
		Transaction 	[]Transaction 	`json:"-"`

	}
)