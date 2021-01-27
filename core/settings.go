package core

import "time"

type Settings struct {
	StaffIds      []string
	DaysDepth     int
	Duration      time.Duration
	TelegramToken string
	ChatID        int64
	Phone         string
	Name          string
	Email         string
}
