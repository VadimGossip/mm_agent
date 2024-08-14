package model

import "time"

type Recipient struct {
	Name     string
	TeamName string
	IsUser   bool
}

type Message struct {
	ID         int64
	Sender     string
	Recipients []Recipient
	Text       string
	SendAt     *time.Time
}
