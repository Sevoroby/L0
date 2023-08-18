# Запустить NATS Streaming Server
    go run NATSStreamingServer/main.go -cid  test -store file -dir store
# Запустить сервис
    go run service/main.go
# Запустить программу для публикации данных в канал
    go run publisher/main.go
# Запустить стресс-тестирование сервиса с помощью Vegeta
    go run vegetaTesting/main.go
# прогнать тесты
	go test ./service
    go test ./publisher