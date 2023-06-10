package model

import (
	"time"
)

type Favorite struct {
	Pk        string    `json:"pk"`
	Sk        string    `json:"sk"`
	UserId    string    `json:"userId"`
	ItemId    string    `json:"itemId"`
	CreatedAt time.Time `json:"created_at"`
}
