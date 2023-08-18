package main

import (
	. "L0/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/nats-io/stan.go"
	"strconv"
	"time"
)

func main() {
	//Подключение к серверу NATS Streaming
	sc, err := stan.Connect("test", "publisher", stan.NatsURL("nats://localhost:4222"))
	//Вывести ошибку
	if err != nil {
		panic(err)
	}
	// Закрыть соединение после окончания
	defer sc.Close()
	for i := 0; i < 100; i++ {
		sendMessage(generateOrder(i), sc)
		//time.Sleep(5 * time.Second)
	}
}
func sendMessage(obj interface{}, sc stan.Conn) bool {
	msgBytes, err := json.Marshal(obj)
	if err != nil {
		fmt.Println("Ошибка преобразования в JSON! " + err.Error())
		return false
	} else {
		fmt.Println(string(msgBytes))
		sc.Publish("message", msgBytes)
		fmt.Println("Сообщение отправлено")
		return true
	}
}

func generateOrder(i int) Order {
	dateCreated, _ := time.Parse("2006-01-02T03:04:05Z", "2021-11-26T06:22:19Z")
	order := Order{
		OrderUid:    "b563feb7b2b84b6test" + strconv.Itoa(i),
		TrackNumber: "WBILMTESTTRACK",
		Entry:       sql.NullString{"WBIL", true},
		Delivery: &Delivery{
			Name:    "Test Testov",
			Phone:   "+9720000000",
			Zip:     "2639809",
			City:    "Kiryat Mozkin",
			Address: "Ploshad Mira 15",
			Region:  "Kraiot",
			Email:   "test@gmail.com",
		},
		Payment: &Payment{
			Transaction:  "b563feb7b2b84b6test0",
			RequestId:    "",
			Currency:     "USD",
			Provider:     "wbpay",
			Amount:       1817,
			PaymentDt:    1637907727,
			Bank:         "alpha",
			DeliveryCost: 1500,
			GoodsTotal:   317,
			CustomFee:    0,
		},
		Items: &[]Item{
			{
				ChrtId:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest00",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmId:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
			{
				ChrtId:      8934930,
				TrackNumber: "WBILMTESTTRACK2",
				Price:       300,
				Rid:         "ac4219087a764ae0btest00",
				Name:        "Thing",
				Sale:        10,
				Size:        "0",
				TotalPrice:  270,
				NmId:        1389212,
				Brand:       "Any",
				Status:      202,
			},
		},
		Locale:            sql.NullString{"en", true},
		InternalSignature: sql.NullString{"", true},
		CustomerId:        sql.NullString{"test", true},
		DeliveryService:   sql.NullString{"meets", true},
		Shardkey:          sql.NullString{"9", true},
		SmId:              sql.NullInt32{9, true},
		DateCreated:       sql.NullTime{dateCreated, true},
		OofShard:          sql.NullString{"1", true},
	}
	if i%2 == 0 {
		order.Items = &[]Item{
			{
				ChrtId:      9934930,
				TrackNumber: "WBILMTESTTRACK",
				Price:       453,
				Rid:         "ab4219087a764ae0btest00",
				Name:        "Mascaras",
				Sale:        30,
				Size:        "0",
				TotalPrice:  317,
				NmId:        2389212,
				Brand:       "Vivienne Sabo",
				Status:      202,
			},
		}
	}
	return order
}
