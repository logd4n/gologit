package rabbitmq

import (
	"encoding/json"
	"log"
	"logger/internal/database"
	"logger/internal/models"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

func ReadMessages() error {
	//1.
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Printf("Не удалось подключиться к RabbitMQ (адрес:\"amqp://guest:guest@rabbimq:5672/\") ")
		log.Printf("Error: %v", err.Error())
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("Не удалось создать канал! Error: %v\n", err.Error())
		return err
	}
	defer ch.Close()

	//2.
	err = ch.ExchangeDeclare(
		"logs_exchange",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Не удалось объявить exchange! Error: %v\n", err.Error())
		return err
	}

	//3.
	queue, err := ch.QueueDeclare(
		"q.database.saver",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Не удалось объявить queue! Error: %v\n", err.Error())
		return err
	}

	//4.
	err = ch.QueueBind(
		queue.Name,
		"logs.#",
		"logs_exchange",
		false,
		nil,
	)
	if err != nil {
		log.Printf("Не удалось связать exchange и queue! Error: %v\n", err.Error())
		return err
	}

	//5.
	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Printf("Не удалось прочитать сообщение из очереди! Error: %v\n", err.Error())
		return err
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		for m := range msgs {
			log.Printf("MESSAGE: %s\n", m.Body)
			m.Ack(false)

			var logMessage models.LogMessage

			err := json.Unmarshal(m.Body, &logMessage)
			if err != nil {
				log.Printf("Не удалось перевести byte[] в json: %s", err.Error())
			}

			err = database.SendLogs(
				&models.LogsTable{
					Level:      string(logMessage.Level),
					Message:    logMessage.Message,
					Created_at: m.Timestamp,
				},
			)
			if err != nil {
				log.Printf(err.Error())
			}
		}
	}()

	wg.Wait()

	log.Printf("\n\nLOGGER IS ENDING!")
	return nil
}
