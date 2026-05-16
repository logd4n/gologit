package main

import (
	"log"
	"logger/database"
	"logger/rabbitmq"
)

func main() {
	err := database.ConnectDB()
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД!: %v", err.Error())
	}

	log.Printf("LOGGER IS RUNNING!\n\n")

	err = rabbitmq.ReadMessages()
	if err != nil {
		log.Printf("Ошибка получения сообщения: %v", err.Error())
	}

	log.Printf("\n\nLOGGER IS ENDING!")
}
