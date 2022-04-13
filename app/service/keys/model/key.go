package model

import "time"

// Key represents Key entity
type Key struct {
	ID        int       `db:"id"`
	Value     string    `db:"value"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
