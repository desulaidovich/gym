package dto

import "time"

type ScheduleResponse struct {
	Error   bool            `json:"error"`
	Message string          `json:"message,omitempty"`
	Data    []*ScheduleItem `json:"data,omitempty"`
}

type ScheduleItem struct {
	ID          string    `json:"id"`
	Date        time.Time `json:"date"`
	Title       string    `json:"title"`
	Club        string    `json:"club"`
	Direction   string    `json:"direction"`
	Description string    `json:"description"`
	Coach       string    `json:"coach"`
	Level       string    `json:"level"`
	Picture     string    `json:"picture"`
}

type ListScheduleQuery struct {
	Title     *string `json:"title,omitempty"     query:"title"`
	Club      *string `json:"club,omitempty"      query:"club"`
	Direction *string `json:"direction,omitempty" query:"direction"`
	Coach     *string `json:"coach,omitempty"     query:"coach"`
}
