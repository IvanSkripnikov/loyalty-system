package models

const ServiceDatabase = "LoyaltySystemService"

const Success = "success"
const Failure = "failure"

type Redis struct {
	Address  string
	Port     string
	Password string
	DB       int
	Stream   string
}
