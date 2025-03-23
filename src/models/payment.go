package models

const TypeDeposit = "deposit"
const TypePayment = "payment"

type Payment struct {
	ID        int     `gorm:"index;type:int" json:"id"`
	UserID    int     `gorm:"index;type:int" json:"userId"`
	Type      string  `gorm:"type:text" json:"type"`
	Amount    float32 `gorm:"type:float" json:"amount"`
	Created   int     `gorm:"index;type:int" json:"created"`
	Status    int8    `gorm:"index;type:tinyint;default:0" json:"status"`
	RequestID string  `gorm:"index;type:string" json:"requestId"`
}
