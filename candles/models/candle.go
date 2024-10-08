package models

import (
	"github.com/google/uuid"
	"time"
)

const (
	TypeMinute CandleType = iota
	TypeHour
	TypeDay
	TypeWeek
)

type CandleType int

type Candle struct {
	ID        uuid.UUID
	StockISIN string
	Open      string
	Close     string
	Time      time.Time
	Type      CandleType
}
