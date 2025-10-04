package handlers

import (
	"context"
	"database/sql"
	"ecommerce/commons"
	"ecommerce/db"
	"ecommerce/templates"
	"errors"
	"time"

	"github.com/gofiber/fiber/v2"
)

func HomeRoute(c *fiber.Ctx) error {
	err := commons.AuthorizeGet(c)
	if err == nil {
		home := templates.Home()
		c.Set("Content-Type", "text/html")
		return home.Render(c.Context(), c.Response().BodyWriter())
	} else {
		authFailed := templates.AuthFailedPage()
		c.Set("Content-Type", "text/html")
		return authFailed.Render(c.Context(), c.Response().BodyWriter())
	}
}

func SinginRoute(c *fiber.Ctx) error {
	signin := templates.SignIn()
	c.Set("Content-Type", "text/html")
	return signin.Render(c.Context(), c.Response().BodyWriter())
}

func SingupRoute(c *fiber.Ctx) error {
	signup := templates.SignUp()
	c.Set("Content-Type", "text/html")
	return signup.Render(c.Context(), c.Response().BodyWriter())
}

func RegisterUser(c *fiber.Ctx) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	ctx := context.Background()
	sqlDb, err := commons.CreateNewDb()
	if err != nil {
		banners := templates.SingupBanner(err)
		return banners.Render(c.Context(), c.Response().BodyWriter())
	}
	queries := db.New(sqlDb)
	_, err = queries.GetUser(ctx, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			hashed_psw, err := commons.HashPassword(password)
			if err != nil {
				banners := templates.SingupBanner(err)
				return banners.Render(c.Context(), c.Response().BodyWriter())
			}
			_, err = queries.CreateUser(ctx, db.CreateUserParams{Username: username, HashedPassword: hashed_psw})
			if err != nil {
				banners := templates.SingupBanner(err)
				return banners.Render(c.Context(), c.Response().BodyWriter())
			} else {
				banners := templates.SingupBanner(nil)
				return banners.Render(c.Context(), c.Response().BodyWriter())
			}
		} else {
			banners := templates.SingupBanner(err)
			return banners.Render(c.Context(), c.Response().BodyWriter())
		}
	} else {
		banners := templates.SingupBanner(errors.New("user already exists"))
		return banners.Render(c.Context(), c.Response().BodyWriter())
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
		if errors.Is(err, sql.ErrNoRows) {
			banners := templates.SingupBanner(errors.New("there is no user with this username"))
			return banners.Render(c.Context(), c.Response().BodyWriter())
		} else {
			banners := templates.SingupBanner(err)
			return banners.Render(c.Context(), c.Response().BodyWriter())
		}
	}
	if !commons.CompareHashToPassword(password, user.HashedPassword) {
		banners := templates.SingupBanner(errors.New("wrong username or password"))
		return banners.Render(c.Context(), c.Response().BodyWriter())
	} else {
		sess_token, errSes := commons.GenerateToken(32)
		csrf_token, errCsrf := commons.GenerateToken(32)
		if errSes != nil || errCsrf != nil {
			banners := templates.SingupBanner(errors.New("an error occurred while generating your authentication credentials"))
			return banners.Render(c.Context(), c.Response().BodyWriter())
		}
		err = queries.UpdateUserTokensLogin(ctx, db.UpdateUserTokensLoginParams{SessionToken: sql.NullString{String: sess_token, Valid: true}, CsrfToken: sql.NullString{String: csrf_token, Valid: true}, Username: username})
		if err != nil {
			banners := templates.SingupBanner(err)
			return banners.Render(c.Context(), c.Response().BodyWriter())
		} else {
			c.Cookie(&fiber.Cookie{
				Name:     "session_token",
				Value:    sess_token,
				Expires:  time.Now().Add(24 * time.Hour),
				HTTPOnly: true,
			})
			c.Cookie(&fiber.Cookie{
				Name:     "csrf_token",
				Value:    csrf_token,
				Expires:  time.Now().Add(24 * time.Hour),
				HTTPOnly: false,
			})
			c.Set("HX-Redirect", "/")
			return c.SendStatus(fiber.StatusOK)
		}
	}
}

func LogoutUser(c *fiber.Ctx) error {
	err := commons.AuthorizePost(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
	} else {
		ctx := context.Background()
		sqlDb, err := commons.CreateNewDb()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Internal server error: " + err.Error()})
		}
		queries := db.New(sqlDb)
		st := c.Cookies("session_token", "")
		csrf := c.Cookies("csrf_token", "")
		queries.UpdateUserTokensLogout(ctx, db.UpdateUserTokensLogoutParams{SessionToken: sql.NullString{String: st, Valid: true}, CsrfToken: sql.NullString{String: csrf, Valid: true}})
		c.Set("HX-Redirect", "/signin")
		return c.SendStatus(fiber.StatusOK)
	}
}
