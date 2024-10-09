package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/chenshuiluke/di-pocket-watcher/api/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5/pgxpool"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	googleOauthConfig *oauth2.Config
)

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/gmail.readonly"},
		Endpoint:     google.Endpoint,
	}
}

func handleGoogleLogin(c *fiber.Ctx) error {
	url := googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

func handleGoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(c.Context(), code)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to exchange token")
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user info")
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to decode user info")
	}

	// TODO: create or update the user in the db and create session or JWT

	return c.JSON(fiber.Map{
		"message": "Successfully authenticated with Google",
		"user":    userInfo,
	})
}

func main() {
	ctx := context.Background()
	url := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME"))

	conn, err := pgxpool.New(ctx, url)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}

	defer conn.Close()

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		queries := db.New(conn)
		results, err := queries.ListUsers(ctx)
		if err != nil {
			log.Println(err)
		}
		return c.JSON(results)
	})

	app.Get("/auth/google/login", handleGoogleLogin)
	app.Get("/auth/google/callback", handleGoogleCallback)

	app.Listen(":8080")
}
