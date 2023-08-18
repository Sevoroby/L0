package model

import (
	"database/sql"
)

type Order struct { //Заказ
	OrderUid          string         `json:"order_uid"`
	TrackNumber       string         `json:"track_number"`
	Entry             sql.NullString `json:"entry"`
	Delivery          *Delivery      `json:"delivery"`
	Payment           *Payment       `json:"payment"`
	Items             *[]Item        `json:"items"`
	Locale            sql.NullString `json:"locale"`
	InternalSignature sql.NullString `json:"internal_signature"`
	CustomerId        sql.NullString `json:"customer_id"`
	DeliveryService   sql.NullString `json:"delivery_service"`
	Shardkey          sql.NullString `json:"shardkey"`
	SmId              sql.NullInt32  `json:"sm_id"`
	DateCreated       sql.NullTime   `json:"date_created"`
	OofShard          sql.NullString `json:"oof_shard"`
}
type Delivery struct { //Доставка
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}
type Item struct { //Товар
	ChrtId      int    `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       int    `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        int    `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  int    `json:"total_price"`
	NmId        int    `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}
type Payment struct { //Оплата
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       int    `json:"amount"`
	PaymentDt    int    `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost int    `json:"delivery_cost"`
	GoodsTotal   int    `json:"goods_total"`
	CustomFee    int    `json:"custom_fee"`
}
