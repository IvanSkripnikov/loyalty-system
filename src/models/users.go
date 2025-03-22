package models

type User struct {
	ID        int    `gorm:"index;type:int" json:"id"`
	UserName  string `gorm:"type:text" json:"username"`
	Password  string `gorm:"type:text" json:"password"`
	FirstName string `gorm:"type:text" json:"first_name"`
	LastName  string `gorm:"type:text" json:"last_name"`
	Email     string `gorm:"type:text" json:"email"`
	Phone     string `gorm:"type:text" json:"phone"`
	Created   int    `gorm:"type:bigint" json:"created"`
	Updated   int    `gorm:"type:bigint" json:"updated"`
	Active    int    `gorm:"index;type:int" json:"active"`
}
