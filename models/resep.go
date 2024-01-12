package models

import "time"

type (
	Recipe struct {
		ID      			uint     		`gorm:"primary_key" json:"id"`
		No      			int       		`json:"no"`
		Tanggal 			time.Time 		`json:"tanggal"`
		Nama_pasien 		string 			`json:"nama_pasien"`
		Nama_dokter 		string 			`json:"nama_dokter"`
		Obat_resep 			string 			`json:"obat_resep"`
		Jumlah_obat_resep 	int 			`json:"jumlah_obat_resep"`
		Flag 				int 			`json:"flag"`
		Transaction 		[]Transaction 	`json:"-"`
	}
)