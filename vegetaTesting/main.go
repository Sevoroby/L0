package main

import (
	"fmt"
	"github.com/nats-io/stan.go"
	vegeta "github.com/tsenart/vegeta/lib"
	"os"
	"time"
)

func main() {
	// Стресс-тест сервиса на получение сообщений NATS Streaming
	attackNATSStreaming()
	// Стресс-тест сервиса на обработку http запросов
	attackHTTP()
}
func attackNATSStreaming() {
	// Подключение к NATS Streaming
	sc, err := stan.Connect("test", "vegeta-stress-test-client", stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		fmt.Println("Ошибка подключения к NATS Streaming: %v", err)
	}

	defer sc.Close()

	// Создаем новый рейт с нагрузкой 200 сообщений в секунду в течение 10 секунд
	rate := vegeta.Rate{Freq: 200, Per: time.Second}
	duration := 10 * time.Second

	// Пустой таргетер
	targeter := vegeta.NewStaticTargeter(vegeta.Target{})

	// Создаем нового атакующего
	attacker := vegeta.NewAttacker()

	// Цикл атаки на сервис сообщениями NATS Streaming
	start := time.Now()
	for range attacker.Attack(targeter, rate, duration, "Stress test") {
		sc.Publish("message", []byte("Test message"))
	}
	//Если количество сообщений в секунду превышает 200, то сервер не справляется с нагрузкой и обрабатывает их дольше 10 сек
	elapsed := time.Since(start)
	fmt.Printf("Время выполнения: %s\n", elapsed)
}
func attackHTTP() {
	// Создаем новый рейт с нагрузкой 4000 запросов в секунду в течение 10 секунд
	rate := vegeta.Rate{Freq: 4000, Per: time.Second}
	duration := 10 * time.Second

	// Создаем новый таргетер с URL локального http-сервера
	targeter := vegeta.NewStaticTargeter(vegeta.Target{
		Method: "GET",
		URL:    "http://localhost:8080/order/b563feb7b2b84b6test0",
	})

	// Создаём новый атакер
	attacker := vegeta.NewAttacker()

	var metrics vegeta.Metrics

	// Атакуем сервис и записываем метрики
	for res := range attacker.Attack(targeter, rate, duration, "Stress test") {
		metrics.Add(res)
	}

	metrics.Close()

	//Если слать более 4000 запросов в секунду, то сервер не способен все обработать и возникает ошибка
	//An operation on a socket could not be performed because the system lacked sufficient buffer space or because a queue was full.
	vegeta.NewTextReporter(&metrics).Report(os.Stdout)
}
