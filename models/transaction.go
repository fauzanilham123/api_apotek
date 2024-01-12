package models

import "time"

type (
	Transaction struct {
		ID        	uint      	`gorm:"primary_key" json:"id"`
		No 			int 		`json:"no"`
		Tanggal 	time.Time 	`json:"tanggal"`
		Nama_kasir 	string 		`json:"nama_kasir"`
		Total_bayar int			`json:"total_bayar"`
		UserID   	uint      	`gorm:"column:user_id" json:"id_user"`
		DrugID   	uint      	`gorm:"column:drug_id" json:"id_drug"`
		RecipeID   	uint      	`gorm:"column:recipe_id" json:"id_recipe"`
		Flag 		int 		`json:"flag"`
		User     	User      	`gorm:"foreignKey:UserID" json:"user"`
		Drug     	Drug      	`gorm:"foreignKey:DrugID" json:"obat"`
		Recipe     	Recipe      `gorm:"foreignKey:RecipeID" json:"resep"`
		CreatedAt 	time.Time 	`json:"created_at"`
		UpdatedAt 	time.Time 	`json:"updated_at"`
	}
)