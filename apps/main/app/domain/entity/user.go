package entity

import "time"

type User struct {
	ID            string
	Name          *string
	Email         *string
	LastLoginAt   *time.Time
	IsDeactivated bool
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (e *User) Created() {
	now := time.Now()
	e.CreatedAt = now
	e.UpdatedAt = now
}

func (e *User) Updated() {
	e.UpdatedAt = time.Now()
}
