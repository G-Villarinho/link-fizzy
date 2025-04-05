package models

import "time"

type LinkVisit struct {
	ID        string
	LinkID    string
	IP        string
	Agent     string
	VisitedAt time.Time
}
