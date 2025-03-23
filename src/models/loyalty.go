package models

const PromocodeTypeStatic = 1
const PromocodeTypePercent = 2

const LoyaltyTypePromocode = 1
const LoyaltyTypeNoOrders = 2
const LoyaltyTypeDiscount1 = 3
const LoyaltyTypeDiscount2 = 4
const LoyaltyTypeDiscount3 = 5
const LoyaltyTypeDiscount4 = 6
const LoyaltyTypeCertificate = 7
const LoyaltyTypeTempDiscount = 8

const TriggerMinimalOrdersSum = "trigger_minimal_orders_sum"
const TriggerFirstLevelOrdersSum = "trigger_first_level_orders_sum"
const TriggerSecondLevelOrdersSum = "trigger_second_level_orders_sum"
const TriggerThirdLevelOrdersSum = "trigger_third_level_orders_sum"
const TriggerSwitchVIPUserCategory = "trigger_switch_vip_user_category"

// managerID = 1 - системный
const ManagerIDSystem = 1

type LoyaltyType struct {
	ID          int    `gorm:"index;type:int" json:"id"`
	Title       string `gorm:"type:text" json:"title"`
	Description string `gorm:"type:text" json:"description"`
	Created     string `gorm:"index;type:string" json:"created"`
	Updated     string `gorm:"index;type:string" json:"updated"`
	Active      int8   `gorm:"index;type:tinyint" json:"active"`
}

func (s LoyaltyType) TableName() string { return "loyalty_type" }

type Loyalty struct {
	ID        int    `gorm:"index;type:int" json:"id"`
	Title     string `gorm:"type:text" json:"title"`
	ManagerID int    `gorm:"index;type:int" json:"managerId"`
	TypeID    int    `gorm:"type:int" json:"typeId"`
	Created   string `gorm:"index;type:string" json:"created"`
	Expired   string `gorm:"index;type:string" json:"expired"`
	Data      string `gorm:"index;type:string" json:"data"`
	Active    int8   `gorm:"index;type:tinyint" json:"active"`
}

func (s Loyalty) TableName() string { return "loyalty" }

type LoyaltyUser struct {
	ID        int  `gorm:"index;type:int" json:"id"`
	UserID    int  `gorm:"index;type:int" json:"userId"`
	LoyaltyID int  `gorm:"index;type:int" json:"loyaltyId"`
	Active    int8 `gorm:"index;type:tinyint" json:"active"`
}

func (s LoyaltyUser) TableName() string { return "loyalty_users" }

type LoyaltyConfiguration struct {
	ID       int    `gorm:"index;type:int" json:"id"`
	Property string `gorm:"index;type:string" json:"property"`
	Value    string `gorm:"index;type:string" json:"value"`
	Active   int8   `gorm:"index;type:tinyint" json:"active"`
}

func (s LoyaltyConfiguration) TableName() string { return "loyalty_configuration" }

// 0 - на все категории товаров, 0 - на все товары
type Promocode struct {
	Type         int8 `json:"type"`
	Value        int  `json:"value"`
	ItemCategory int8 `json:"itemCategory"`
	Item         int8 `json:"item"`
}

type Discount struct {
	Value        int  `json:"value"`
	ItemCategory int8 `json:"itemCategory"`
	Item         int8 `json:"item"`
}

type Certificate struct {
	Value        int  `json:"value"`
	ItemCategory int8 `json:"itemCategory"`
	Item         int8 `json:"item"`
}

type TempDiscount struct {
	Type         int8   `json:"type"`
	Value        int    `json:"value"`
	ItemCategory int8   `json:"itemCategory"`
	Item         int8   `json:"item"`
	FromDate     string `json:"fromDate"`
	ToDate       string `json:"toDate"`
}

type FirstDiscount struct {
	Type         int8 `json:"type"`
	Value        int  `json:"value"`
	ItemCategory int8 `json:"itemCategory"`
	Item         int8 `json:"item"`
}
