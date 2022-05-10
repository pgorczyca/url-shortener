package model

import (
	"time"
)

type Url struct {
	Long      string    `json:"long"`
	Short     string    `json:"short"`
	ExpiredAt time.Time `json:"exipred_at"`
	CreatedAt time.Time `json:"created_at"`
}
