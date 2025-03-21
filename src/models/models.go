package models

const ServiceDatabase = "LoyaltySystemService"

type Redis struct {
	Address  string
	Port     string
	Password string
	DB       int
	Stream   string
}
