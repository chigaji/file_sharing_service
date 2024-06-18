package models

import "time"

type FileMetadata struct {
	Filename string    `db:"filename"`
	Expiry   time.Time `db:"expiry"`
}

type User struct {
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
}
