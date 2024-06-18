package storage

import (
	"time"

	"github.com/chigaji/file_sharing_service/pkg/models"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var db *sqlx.DB

func InitDB() error {

	var err error

	// db, err = sqlx.Connect("postgres", "user=admin dbname=file_sharing_db sslmode=disable")
	db, err = sqlx.Connect("postgres", "postgres://admin:admin1!@localhost:5432/file_sharing_db?sslmode=disable")

	if err != nil {
		return err
	}

	schema := `
	CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password TEXT NOT NULL
	);

	CREATE TABLE IF NOT EXISTS files (
		id UUID PRIMARY KEY, 
		filename TEXT NOT NULL,
		expiry TIMESTAMP NOT NULL
	);
	`

	_, err = db.Exec(schema)
	if err != nil {
		return err
	}

	return nil
}
func SaveUser(user models.User) error {
	_, err := db.Exec("INSERT INTO users (username, password) VALUES ($1, $2)", user.Username, user.Password)
	return err
}
func SaveFileData(fileID string, filename string, expiry time.Time) error {

	_, err := db.Exec("INSERT INTO file (id, filename, expiry) VALUES ($1, $2, $3)", fileID, filename, expiry)
	if err != nil {
		return err
	}
	return nil
}

func GetFileData(fileID string) (models.FileMetadata, error) {

	var fileMetadata models.FileMetadata

	err := db.Get(&fileMetadata, "SELECT * FROM files WHERE id = $1", fileID)

	// if err != nil {
	// 	return fileMetadata, err
	// }

	return fileMetadata, err

}

func GetUserPassword(username string) (string, error) {
	var userPassword string

	err := db.Get(&userPassword, "SELECT password FROM users WHERE username = $1", username)
	return userPassword, err
}
