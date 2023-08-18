package main

import (
	"github.com/nats-io/stan.go"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Тест для проверки заполнения кэша из бд
func TestFillCacheMoreThanZero(t *testing.T) {
	fillCache()
	assert.Equal(t, true, len(cache) > 0)
}

// Тест на пустой объект во входящем сообщении
func TestGetMessageHandlerEmpty(t *testing.T) {
	m := new(stan.Msg)
	m.Data = []byte("{}")
	getMessageHandler(m)
}

// Тест на неверный объект во входящем сообщении
func TestGetMessageHandlerWrongObject(t *testing.T) {
	m := new(stan.Msg)
	m.Data = []byte("{\"Make\":\"toyota\",\"Model\":\"camry\"}\n")
	getMessageHandler(m)
}
