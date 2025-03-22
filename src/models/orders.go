package models

type Order struct {
	ID        int     `json:"id"`
	UserID    int     `json:"userId"`
	ItemID    int     `json:"itemId"`
	Volume    int     `json:"volume"`
	Price     float32 `json:"price"`
	Created   int     `json:"created"`
	Updated   int     `json:"updated"`
	Status    uint8   `json:"status"`
	RequestID string  `json:"requestId"`
}
