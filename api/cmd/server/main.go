package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	connection_manager "github.com/chenshuiluke/di-pocket-watcher/api/internal"
	"github.com/chenshuiluke/di-pocket-watcher/api/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	jwtware "github.com/gofiber/contrib/jwt"
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
	email := userInfo["email"].(string)
	var user db.User
	existingUser, err := connection_manager.Mgr.Queries.GetUserByEmail(c.Context(), email)
	if err == nil {
		user = db.User{
			ID:    existingUser.ID,
			Email: existingUser.Email,
		}
	}
	if (db.User{}) == user {
		params := db.CreateUserParams{
			Email: email,
			Token: token.RefreshToken,
		}
		user, err = connection_manager.Mgr.Queries.CreateUser(c.Context(), params)
		if err != nil {
			log.Println(err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	// Create the Claims
	claims := jwt.MapClaims{
		"id":  user.ID,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}

	// Create token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	t, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	return c.JSON(fiber.Map{"token": t})
}

func main() {

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {

		results, err := connection_manager.Mgr.Queries.ListUsers(c.Context())
		if err != nil {
			log.Println(err)
		}
		return c.JSON(results)
	})

	app.Get("/auth/google/login", handleGoogleLogin)
	app.Get("/auth/google/callback", handleGoogleCallback)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
	}))

	app.Listen(":8080")
}
