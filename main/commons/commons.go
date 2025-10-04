package commons

import (
	"crypto/rand"
	"ecommerce/db"
	"encoding/base64"
	"errors"

	"context"
	"database/sql"
	_ "embed"

	"github.com/gofiber/fiber/v2"
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

var ErrUnauthorized = errors.New("unauthorized")

func AuthorizePost(c *fiber.Ctx) error {
	sqlDb, err := CreateNewDb()
	if err != nil {
		return ErrUnauthorized
	}
	st := c.Cookies("session_token", "")
	if st == "" {
		return ErrUnauthorized
	}
	queries := db.New(sqlDb)
	ctx := context.Background()
	user, err := queries.GetUserBySessionToken(ctx, sql.NullString{String: st, Valid: true})
	if err != nil {
		return ErrUnauthorized
	}
	csrf := c.Cookies("csrf_token", "")
	if csrf == "" {
		return ErrUnauthorized
	}
	if csrf != user.CsrfToken.String {
		return ErrUnauthorized
	}
	return nil
}

func AuthorizeGet(c *fiber.Ctx) error {
	sqlDb, err := CreateNewDb()
	if err != nil {
		return ErrUnauthorized
	}
	st := c.Cookies("session_token", "")
	if st == "" {
		return ErrUnauthorized
	}
	queries := db.New(sqlDb)
	ctx := context.Background()
	_, err = queries.GetUserBySessionToken(ctx, sql.NullString{String: st, Valid: true})
	if err != nil {
		return ErrUnauthorized
	}
	return nil
}
