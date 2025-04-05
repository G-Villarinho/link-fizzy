package models

import "time"

type Session struct {
	ID       string
	UserID   string
	Token    string
	IP       string
	Agent    string
	CreateAt time.Time
	ExpireAt time.Time
}
