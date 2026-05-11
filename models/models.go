package models

import "time"

type LogsTable struct {
	Id         uint32    `json:"id"`
	Level      string    `json:"level"`
	Message    string    `json:"message"`
	Created_at time.Time `json:"created_at"`
}

type LogLevel string

const (
	Error LogLevel = "error"
	Info  LogLevel = "info"
	Panic LogLevel = "panic"
	Warn  LogLevel = "warning"
	Fatal LogLevel = "fatal"
	Debug LogLevel = "debug"
)

type LogMessage struct {
	Level   LogLevel `json:"level"`
	Message string   `json:"message"`
}
