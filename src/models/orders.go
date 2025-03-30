package models

type Order struct {
	ID          int     `json:"id"`
	UserID      int     `json:"userId"`
	ItemID      int     `json:"itemId"`
	Volume      int     `json:"volume"`
	Price       float32 `json:"price"`
	Created     string  `json:"created"`
	Updated     string  `json:"updated"`
	Status      int     `json:"status"`
	RequestID   string  `json:"requestId"`
	Promocode   string  `json:"promocode"`
	Certificate string  `json:"certificate"`
}
