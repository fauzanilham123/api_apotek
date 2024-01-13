package models

import "time"

type (
	Recipe struct {
		ID      			uint     		`gorm:"primary_key" json:"id"`
		No      			int       		`gorm:"unique" json:"no"`
		Tanggal 			string 			`json:"tanggal"`
		Nama_pasien 		string 			`json:"nama_pasien"`
		Nama_dokter 		string 			`json:"nama_dokter"`
		Nama_obat 			string 			`json:"obat_resep"`
		Jumlah_obat_resep 	int 			`json:"jumlah_obat_resep"`
		Flag 				int 			`json:"flag"`
		CreatedAt 			time.Time 		`json:"created_at"`
		UpdatedAt 			time.Time 		`json:"updated_at"`
		Transaction 		[]Transaction 	`json:"-"`
	}
)