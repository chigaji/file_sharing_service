package models

import "time"

type FileMetadata struct {
	Filename string    `db:"filename"`
	Expiry   time.Time `db:"expiry"`
}

type User struct {
	Username string `db:"username"`
	Password string `db:"password"`
}
