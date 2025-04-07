package types

import (
	"fmt"
	"time"
)

type Weekday int

const (
	Monday Weekday = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

type TimeRange struct {
	Opening time.Time `json:"opening"`
	Closing time.Time `json:"closing"`
}

type WeeklySchedule struct {
	Hours map[Weekday][]TimeRange `json:"hours"`
}

func NewWeeklySchedule() WeeklySchedule {
	return WeeklySchedule{
		Hours: make(map[Weekday][]TimeRange),
	}
}

func (ws *WeeklySchedule) AddOpeningHours(day Weekday, opening, closing time.Time) error {
	if closing.Before(opening) {
		return fmt.Errorf("closing time must be after opening time")
	}

	ws.Hours[day] = append(ws.Hours[day], TimeRange{
		Opening: opening,
		Closing: closing,
	})
	return nil
}