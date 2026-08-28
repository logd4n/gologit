## GOLOGIT - готовая система записи и хранения логов в БД.

Если Вашему проекту требуется готовое и простое решение для записи и хранения логов, то ***Gologit*** является классным вариантом!

***Gologit*** использует RabbitMQ для приема сообщений (логов) от Вашего сервиса и PostgreSQL для их хранения.

### <u>Для работы системы требуется</u>:
- docker-compose
- RabbitMQ
- PostgreSQL *(можно использовать базу Вашего сервиса)*
- Настроить Healthcheck для RabbitMQ *(крайне желательно)*

### Принцип работы логгера:
***Gologit*** подключается к БД и проверяет наличие таблицы *"logs"*, если ее нет - логгер отправит запрос на создание. 
Далее ***Gologit*** подключается к *RabbitMQ* (драйвер amqp) и берет сообщения из очереди *"q.database.saver"* (сообщения распределяет exchange *"logs_exchange"*).
После получения сообщения логгер отправляет эти данные в БД.

***Gologit*** имеет админ-панель для просмотра логов, которая доступна <u>только локальном хосте</u>, то есть для ее запуска в брауезере Вы обращаетесь к localhost:*\*port\**. Порт выделяете самостоятельно, так как в Вашем случае он может быть занят другим процессом. Внутри контейнера сервер логгера слушает порт 8081.

### В каком формате должны быть данные:
***Gologit*** принимает данные в виде структуры, которую потом парсит и отправляет в БД.

#### Форматы данных:
#### Структура сообщения:
``` golang
type LogMessage struct {
	Level   LogLevel `json:"level"`
	Message string   `json:"message"`
}
```

#### Таблица logs:
``` golang
type LogsTable struct {
	Id         uint32    `json:"id"`
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	Created_at time.Time `json:"created_at"`
}
```

#### Типы ошибок:
``` golang
type LogLevel string

const (
	Error LogLevel = "error"
	Info  LogLevel = "info"
	Panic LogLevel = "panic"
	Warn  LogLevel = "warning"
	Fatal LogLevel = "fatal"
	Debug LogLevel = "debug"
)
```

### Примеры использования:
В качестве примеров будет использован язык Go.
#### Функция вывода и передачи сообщения брокеру:
``` golang
var SendMsgErr = errors.New("Не удалось отправить сообщение брокеру!")

func LogPrint(message string, level models.LogLevel) error {
	log.Printf("%s\n", message)

	err := NewMessage(models.LogMessage{
		Level:   level,
		Message: message,
	})
	if err != nil {
		return SendMsgErr
	}

	return nil
}
```

#### Кастомные ошибки:
``` golang
var (
	//Errors
	ConnErr         = errors.New("Ошибка подключения")
	ConnTmtErr      = errors.New("Превышено время ожидания RabbitMQ!")
	DeclareChErr    = errors.New("Не удалось создать канал!")
	DeclareExErr    = errors.New("Не удалось создать Exchange!")
	DeclareQueueErr = errors.New("Не удалось создать очередь!")
	QueueBindErr    = errors.New("Не удалось связать exchange и queue!")
	MarshallErr     = errors.New("Ошибка перевода структуры в byte[]: ")
	PublishErr      = errors.New("Не удалось опубликовать сообщение!")
)
```

#### Функция подключения к брокеру:
``` golang
func ConnectionAttempt() error {
	var conn *amqp.Connection
	var err error

	for i := 1; i <= 10; i++ {
		log.Printf("Попытка подключения к RabbitMQ №%d...\t\n", i)

		conn, err = amqp.Dial("amqp://guest:guest@rabbitmq:5672/")

		if err == nil {
			conn.Close()
			break
		}

		log.Printf("%v: \"%v\"\n", ConnErr.Error(), err.Error())
		time.Sleep(5 * time.Second)
	}

	if err != nil {
		log.Printf("%v Ошибка: %v\n", ConnTmtErr, err.Error())
		return err
	}

	log.Printf("Подключение к RabbitMQ выполнено успешно!")
	return nil
}
```
#### Функция отправки сообщения:
``` golang
func NewMessage(message models.LogMessage) error {
	// 1.
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq:5672/")
	if err != nil {
		log.Printf("%v (адрес:\"amqp://guest:guest@rabbitmq:5672/\"): %v\n",
			ConnErr.Error(), err.Error())
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Printf("%v: %v\n", DeclareChErr.Error(), err.Error())
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
		log.Printf("%v: %v\n", DeclareExErr.Error(), err.Error())
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
		log.Printf("%v: %v\n", DeclareQueueErr.Error(), err.Error())
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
		log.Printf("%v: %v\n", QueueBindErr.Error(), err.Error())
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//5.
	body, err := json.Marshal(message)
	if err != nil {
		return MarshallErr
	}

	err = ch.PublishWithContext(ctx, "logs_exchange",
		"logs."+string(message.Level),
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
			Timestamp:   time.Now().UTC(),
		},
	)
	if err != nil {
		log.Printf("%v: %v\n", PublishErr.Error(), err.Error())
		return err
	}

	return nil
}
```

#### Таким образом Вы можете использовать эту простую систему в своих проектах :)

## ***GO LOG IT!***