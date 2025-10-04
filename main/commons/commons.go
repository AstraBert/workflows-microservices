package commons

import (
	"crypto/rand"
	"encoding/base64"

	"context"
	"database/sql"
	_ "embed"

	"golang.org/x/crypto/bcrypt"

	_ "modernc.org/sqlite"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CompareHashToPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func GenerateToken(tokenLength int) (string, error) {
	bytes := make([]byte, tokenLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	} else {
		return base64.URLEncoding.EncodeToString(bytes), nil
	}
}

//go:embed schema.sql
var ddl string

func CreateNewDb() (*sql.DB, error) {
	ctx := context.Background()

	db, err := sql.Open("sqlite", "users.db")
	if err != nil {
		return nil, err
	}

	// create tables
	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	return db, nil
}
