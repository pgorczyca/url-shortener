package http

import "time"

type CreateUrlRequest struct {
	Long string `json:"long" validate:"required,url"`
}
type UrlResponse struct {
	Long      string    `json:"long"`
	Short     string    `json:"short"`
	CreatedAt time.Time `json:"created_at"`
	ExpiredAt time.Time `json:"expired_at"`
}
