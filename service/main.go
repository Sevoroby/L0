package main

import (
	. "L0/model"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/nats-io/stan.go"
	"log"
	"net/http"
	"sync"
)

// Создание мапы в качестве кэша для заказов
var cache = make(map[string]Order)

// Лок для мапы
var lock = sync.RWMutex{}

func main() {
	// Подключение к серверу NATS Streaming
	sc, err := stan.Connect("test", "service", stan.NatsURL("nats://localhost:4222"))
	// Паника, если ошибка подключения
	if err != nil {
		panic(err)
	}
	// Закрыть соединение после окончания
	defer sc.Close()
	// Подписка на тему message и задание обработчика getMessageHandler,
	//а также задание DurableName, чтобы восстанавливать сообщения из очереди во время падения сервиса
	go sc.Subscribe("message", getMessageHandler, stan.DurableName("service"))
	// Создаём и запускаем http-сервер на порту 8080
	go createServer()

	//Заполнение кэша из БД
	go fillCache()

	//Ожидание
	w := sync.WaitGroup{}
	w.Add(1)
	w.Wait()
}

func createServer() {
	// Используем gin для создания и использования http сервера
	router := gin.Default()
	// Установить папку templates для шаблонов HTML
	router.LoadHTMLGlob("service/templates/*.html")
	// Установить маршрут для GET запросов по определённмоу id заказа
	router.GET("/order/:id", getOrderByIdHandler)
	// Запуск сервера на порту 8080
	router.Run(":8080")
}
func fillCache() int {
	//Подключение к базе данных
	db := connectToDB()
	//Закрыть соединение после окончания
	defer db.Close()
	// Выполнение SELECT запроса
	query := "SELECT * FROM public.\"Order\""
	rows, err := db.Query(query)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	// Чтение результатов запроса
	for rows.Next() {
		o := Order{}
		var deliveryBytes []uint8
		var paymentBytes []uint8
		var itemsBytes []uint8
		err := rows.Scan(&o.OrderUid, &o.TrackNumber, &o.Entry, &o.Locale, &o.InternalSignature, &o.CustomerId, &o.DeliveryService, &o.Shardkey, &o.SmId, &o.DateCreated, &o.OofShard, &deliveryBytes, &itemsBytes, &paymentBytes)
		if err != nil {
			log.Fatal(err)
		}
		if deliveryBytes != nil {
			err = json.Unmarshal(deliveryBytes, &o.Delivery)
			if err != nil {
				fmt.Println("Ошибка преобразования из JSON в объект Delivery")
			}
		}
		if itemsBytes != nil {
			err = json.Unmarshal(itemsBytes, &o.Items)
			if err != nil {
				fmt.Println("Ошибка преобразования из JSON в объект Items")
			}
		}
		if paymentBytes != nil {
			err = json.Unmarshal(paymentBytes, &o.Payment)
			if err != nil {
				fmt.Println("Ошибка преобразования из JSON в объект Payment")
			}
		}
		// Если в мапе нет заказа с таким id, то добавляем его
		lock.RLock()
		_, inMap := cache[o.OrderUid]
		lock.RUnlock()
		if !inMap {
			lock.Lock()
			cache[o.OrderUid] = o
			lock.Unlock()
		}
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}
	return len(cache)
}
func getMessageHandler(m *stan.Msg) {
	var order Order
	fmt.Println("Получено сообщение")
	fmt.Println(string(m.Data))
	err := json.Unmarshal(m.Data, &order)
	// Если пришёл объект, отличающийся от Order, то выводим ошибку и пропускаем
	if err != nil || string(m.Data) == "{}" || order.OrderUid == "" {
		fmt.Println("Неверный объект")
	} else {

		// Проверка, есть ли в кэше заказ с таким id
		lock.RLock()
		_, inMap := cache[order.OrderUid]
		lock.RUnlock()
		if !inMap {
			lock.Lock()
			cache[order.OrderUid] = order
			lock.Unlock()
		}
		// Добавление записи в БД
		db := connectToDB()
		insertToDB(db, order)
	}
}

func connectToDB() *sql.DB {
	driverName := "postgres"
	userName := "SVorobiev"
	pass := "m5t856415lekord"
	addr := "localhost"
	dbName := "test"
	// Открытие соединения с бд
	db, err := sql.Open(driverName, "postgresql://"+userName+":"+pass+"@"+addr+"/"+dbName+"?sslmode=disable")
	// Выводим ошибку, если бд недоступна
	if err != nil {
		fmt.Println(err.Error())
	}
	return db
}
func insertToDB(db *sql.DB, order Order) {
	//Подключение к базе данных
	db = connectToDB()
	defer db.Close()
	//Закрыть соединение после окончания
	// Преобразование в JSON объектов Delivery, Payment и Items
	jsonDelivery, err := json.Marshal(order.Delivery)
	if err != nil {
		panic(err)
	}
	jsonPayment, err := json.Marshal(order.Payment)
	if err != nil {
		panic(err)
	}
	jsonItems, err := json.Marshal(order.Items)
	if err != nil {
		panic(err)
	}
	// Добавление записи
	query := "INSERT INTO public.\"Order\" (\"OrderUid\",\"TrackNumber\",\"Entry\",\"Delivery\",\"Payment\",\"Items\",\"Locale\",\"InternalSignature\",\"CustomerId\",\"DeliveryService\",\"Shardkey\",\"SmId\",\"DateCreated\",\"OofShard\") VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14)"
	_, err = db.Exec(query, order.OrderUid, order.TrackNumber, order.Entry, jsonDelivery, jsonPayment, jsonItems, order.Locale, order.InternalSignature, order.CustomerId, order.DeliveryService, order.Shardkey, order.SmId, order.DateCreated, order.OofShard)
	// Если запись не была добавлена, то выводим ошибку и пропускаем
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Запись успешно добавлена в бд")
	}
}
func getOrderByIdHandler(c *gin.Context) {
	// Получение значения параметра "id" из URL
	id := c.Param("id")

	var order Order
	// Поиск в мапе заказа по id
	lock.RLock()
	val, inMap := cache[id]
	lock.RUnlock()
	if inMap {
		order = val
	} else {
		c.String(http.StatusNotFound, "Заказ с таким id не найден")
		return
	}

	// Шаблон HTML для заказов
	c.HTML(http.StatusOK, "order.html", gin.H{"Order": order})
}
