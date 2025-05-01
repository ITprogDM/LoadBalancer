package models

type Client struct {
	ID         string `json:"client_id"`
	Capacity   int    `json:"capacity"`
	RatePerSec int    `json:"rate_per_sec"`
}
