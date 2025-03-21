package models

type Loyalty struct {
	ID          int    `gorm:"index;type:int" json:"id"`
	UserID      int    `gorm:"index;type:int" json:"userId"`
	Type        string `gorm:"type:text" json:"type"`
	Description string `gorm:"type:text" json:"description"`
	Created     int    `gorm:"index;type:bigint" json:"created"`
	Active      int8   `gorm:"index;type:tinyint" json:"active"`
}

func (s Loyalty) TableName() string { return "loyalty" }
