package main

import (
	"github.com/nats-io/stan.go"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Тест с корректным сообщением при запущенном NATS Streaming Server
func TestSendMessageCorrect(t *testing.T) {
	sc, err := stan.Connect("test", "publisher", stan.NatsURL("nats://localhost:4222"))
	//Вывести ошибку
	if err != nil {
		panic(err)
	}
	defer sc.Close()
	// Создание тестируемого объекта
	msg := struct {
		Make  string
		Model string
	}{
		"toyota",
		"camry",
	}

	// Вызов функции sendMessage
	res := sendMessage(msg, sc)

	assert.Equal(t, true, res)
}

// Тест с некорректным сообщением при запущенном NATS Streaming Server
func TestSendMessageIncorrect(t *testing.T) {
	sc, err := stan.Connect("test", "publisher", stan.NatsURL("nats://localhost:4222"))
	//Вывести ошибку
	if err != nil {
		panic(err)
	}
	defer sc.Close()
	res := sendMessage(make(chan int), sc)
	assert.Equal(t, false, res)
}
