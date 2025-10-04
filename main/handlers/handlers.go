package handlers

import (
	"context"
	"database/sql"
	"ecommerce/commons"
	"ecommerce/db"
	"ecommerce/templates"
	"errors"

	"github.com/gofiber/fiber/v2"
)

func HomeRoute(c *fiber.Ctx) error {
	home := templates.Home()
	c.Set("Content-Type", "text/html")
	return home.Render(c.Context(), c.Response().BodyWriter())
}

func RegisterUser(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	ctx := context.Background()
	sqlDb, err := commons.CreateNewDb()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
	}
	queries := db.New(sqlDb)
	_, err = queries.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			hashed_psw, err := commons.HashPassword(password)
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
			}
			_, err = queries.CreateUser(ctx, db.CreateUserParams{Username: username, HashedPassword: hashed_psw})
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
			} else {
				return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User successfully created"})
			}
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
		}
	} else {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"message": "User already exists"})
	}
}

func LoginUser(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	ctx := context.Background()
	sqlDb, err := commons.CreateNewDb()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
	}
	queries := db.New(sqlDb)
	user, err := queries.GetUser(ctx, username)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
	}
	hashed_psw, err := commons.HashPassword(password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
	}
	if user.HashedPassword != hashed_psw {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"message": "Username or Password are wrong"})
	} else {
		sess_token, errSes := commons.GenerateToken(32)
		csrf_token, errCsrf := commons.GenerateToken(32)
		if errSes != nil || errCsrf != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "An error occurred while generating your authentication credentials"})
		}
		err = queries.UpdateUserTokensLogin(ctx, db.UpdateUserTokensLoginParams{SessionToken: sql.NullString{String: sess_token, Valid: true}, CsrfToken: sql.NullString{String: csrf_token, Valid: true}, Username: username})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
		} else {
			return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "User logged in successfully"})
		}
	}
}

func LogoutUser(c *fiber.Ctx) error {
	return nil
}
