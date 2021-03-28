package main

import "time"

type TableUser struct {
	Phone         string    `json:"phone" gorm:"primaryKey"`
	RealName      string    `json:"real_name"`
	UserID        int       `json:"-"`
	UUID          string    `json:"-"`
	CID           string    `json:"-" gorm:"column:cid"`
	Device        string    `json:"-"`
	SessionID     string    `json:"-"`
	StoreID       string    `json:"store_id"`
	Token         string    `json:"-"`
	TicketName    string    `json:"-"`
	IDCardNum     string    `json:"-"`
	Qualification bool      `json:"qualification" gorm:"default:false"`
	RemainCount   int       `json:"remain_count" gorm:"default:0"`
	CreateTime    time.Time `json:"create_time" gorm:"default:(datetime('now', 'localtime'))"`

	TempID string `json:"-" gorm:"-"`
}

func (*TableUser) TableName() string {
	return "user"
}

type TableOrder struct {
	ID        int       `json:"id" gorm:"PRIMARY_KEY;AUTO_INCREMENT" `
	Phone     string    `json:"phone"`
	Num       int       `json:"num"`
	StoreID   int       `json:"store_id"`
	OrderTime time.Time `json:"order_time" gorm:"default:(datetime('now', 'localtime'))"`
}

func (*TableOrder) TableName() string {
	return "order"
}

type TableStockRecord struct {
	StoreID   int       `json:"store_id"`
	StoreName string    `json:"store_name"`
	Stock     int       `json:"stock"`
	QueryTime time.Time `json:"query_time" gorm:"default:(datetime('now', 'localtime'))"`
}

func (*TableStockRecord) TableName() string {
	return "stock_record"
}
