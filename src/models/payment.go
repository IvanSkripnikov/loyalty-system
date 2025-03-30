package models

const TypeDeposit = "deposit"
const TypePayment = "payment"

type Payment struct {
	ID        int     `json:"id"`
	UserID    int     `json:"userId"`
	Type      string  `json:"type"`
	Amount    float32 `json:"amount"`
	Created   string  `json:"created"`
	Status    int     `json:"status"`
	RequestID string  `json:"requestId"`
}
