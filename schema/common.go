package schema

import "github.com/jinzhu/gorm"

type Mappings struct {
	gorm.Model
	URL   string `json:"url" gorm:"column:url"`
	Hash  string `json:"hash" gorm:"column:hash"`
	Count int    `json:"counter" gorm:"column:count"`
}

type Counter struct {
	gorm.Model
	MachineID int `json:"machine_id" gorm:"column:machine_id"`
	Count     int `json:"count" gorm:"column:count"`
}
