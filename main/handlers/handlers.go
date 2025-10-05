package handlers

import (
	"context"
	"database/sql"
	"ecommerce/commons"
	"ecommerce/db"
	"ecommerce/models"
	"ecommerce/templates"
	"encoding/json"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"
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

func ShowOrderPage(c *fiber.Ctx) error {
	err := commons.AuthorizeGet(c)
	if err != nil {
		authFailed := templates.AuthFailedPage()
		c.Set("Content-Type", "text/html")
		return authFailed.Render(c.Context(), c.Response().BodyWriter())
	}
	c.Set("Content-Type", "text/html")
	itemId := c.Params("id", "1")
	val, err := strconv.Atoi(itemId)
	if err != nil || val > 3 || val < 1 {
		page404 := templates.Page404("the id of the item you would like to order should be a number between 1 and 3")
		return page404.Render(c.Context(), c.Response().BodyWriter())
	}
	orderPage := templates.OrderForm(val)
	return orderPage.Render(c.Context(), c.Response().BodyWriter())
}

func HandleOrder(c *fiber.Ctx) error {
	err := commons.AuthorizePost(c)
	if err != nil {
		banners := templates.SingupBanner(errors.New("you are not authorized to place orders"))
		return banners.Render(c.Context(), c.Response().BodyWriter())
	}
	sqlDb, err := commons.CreateNewDb()
	if err != nil {
		banners := templates.SingupBanner(errors.New("you are not authorized to place orders"))
		return banners.Render(c.Context(), c.Response().BodyWriter())
	}
	queries := db.New(sqlDb)
	user, err := queries.GetUserBySessionToken(context.Background(), sql.NullString{String: c.Cookies("session_token"), Valid: true})
	if err != nil {
		banners := templates.SingupBanner(errors.New("you are not authorized to place orders"))
		return banners.Render(c.Context(), c.Response().BodyWriter())
	}
	userId := user.ID
	firstName := c.FormValue("firstName")
	lastName := c.FormValue("lastName")
	address := c.FormValue("address")
	address2 := c.FormValue("address2")
	city := c.FormValue("city")
	state := c.FormValue("state")
	zip := c.FormValue("zip")
	country := c.FormValue("country")
	email := c.FormValue("email")
	phone := c.FormValue("phone")
	paymentMethod := c.FormValue("paymentMethod")
	amount := c.FormValue("amount")
	orderId, err := commons.GenerateToken(16)
	c.Set("Content-Type", "text/html")
	if err != nil {
		banners := templates.SingupBanner(err)
		return banners.Render(c.Context(), c.Response().BodyWriter())
	}
	order := models.Order{FirstName: firstName, LastName: lastName, Address: address, Address2: address2, City: city, State: state, Zip: zip, Country: country, Email: email, Phone: phone, PaymentMethod: paymentMethod, Amount: amount, OrderId: orderId, UserId: userId}
	byteData, err := json.Marshal(order)
	if err != nil {
		banners := templates.SingupBanner(err)
		return banners.Render(c.Context(), c.Response().BodyWriter())
	}
	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "orders", 0)
	if err != nil {
		banners := templates.SingupBanner(err)
		return banners.Render(c.Context(), c.Response().BodyWriter())
	}
	conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	_, err = conn.WriteMessages(
		kafka.Message{Value: byteData},
	)
	if err != nil {
		banners := templates.SingupBanner(err)
		return banners.Render(c.Context(), c.Response().BodyWriter())
	}
	if err := conn.Close(); err != nil {
		banners := templates.SingupBanner(err)
		return banners.Render(c.Context(), c.Response().BodyWriter())
	}
	connStatus, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", "order-status", 0)
	if err != nil {
		banners := templates.SingupBanner(err)
		return banners.Render(c.Context(), c.Response().BodyWriter())
	}
	connStatus.SetReadDeadline(time.Now().Add(10 * time.Minute))
	batch := connStatus.ReadBatch(1e3, 1e6) // fetch 1KB min, 1MB max

	toMonitor := make([]models.OrdersReportStatus, 0, 3) // Changed to hold parsed structs
	b := make([]byte, 1e3)                               // 1KB max per message

	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}

		// Parse the JSON message
		var status models.OrdersReportStatus
		if err := json.Unmarshal(b[:n], &status); err != nil {
			// Log the error or handle it appropriately
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		toMonitor = append(toMonitor, status)
		if len(toMonitor) == 3 {
			break
		}
	}

	if err := batch.Close(); err != nil {
		banners := templates.SingupBanner(err)
		return banners.Render(c.Context(), c.Response().BodyWriter())
	}

	if err := connStatus.Close(); err != nil {
		banners := templates.SingupBanner(err)
		return banners.Render(c.Context(), c.Response().BodyWriter())
	}
	successPayment := false
	successOrder := false
	successStock := false
	for _, item := range toMonitor {
		switch item.Type {
		case "payment":
			successPayment = item.Success
		case "order":
			successOrder = item.Success
		case "stock":
			successStock = item.Success
		}
	}
	orderBanner := templates.OrderBanner(successPayment, successOrder, successStock)
	return orderBanner.Render(c.Context(), c.Response().BodyWriter())
}
