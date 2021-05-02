package types

import "time"

type Schedule struct {
	Id         string    `json:"id"`
	ItemId     string    `json:"item_id"`
	ItemType   string    `json:"item_type"`
	Action     string    `json:"action"`
	CreatedAt  time.Time `json:"created_at"`
	ExecuteOn  time.Time `json:"execute_on"`
	RetryCount int       `json:"retry_count"`
}
