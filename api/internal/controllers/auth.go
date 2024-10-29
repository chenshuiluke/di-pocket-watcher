package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2/log"

	connection_manager "github.com/chenshuiluke/di-pocket-watcher/api/internal"
	"github.com/chenshuiluke/di-pocket-watcher/api/internal/db"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthController struct {
}

type JWTClaims struct {
	ID json.Number `json:"id"`
	jwt.RegisteredClaims
}

func (c *JWTClaims) GetID() (int64, error) {
	return c.ID.Int64()
}

var (
	googleOauthConfig *oauth2.Config
)

func init() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/api/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/gmail.readonly",
		},
		Endpoint: google.Endpoint,
	}
}

func (AuthController) GetCurrentUser(c *fiber.Ctx) error {
	jwtToken := c.Locals("user").(*jwt.Token)
	claims, ok := jwtToken.Claims.(*JWTClaims)
	if !ok {
		log.Error("Failed to parse JWT claims")
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	id, err := claims.GetID()
	if err != nil {
		log.Error("Failed to parse user ID:", err)
		return c.SendStatus(fiber.StatusUnauthorized)
	}

	user, err := connection_manager.Mgr.Queries.GetUser(c.Context(), id)
	if err != nil {
		log.Error("Failed to get user:", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	return c.JSON(fiber.Map{"user": user})
}

func (AuthController) HandleGoogleLogin(c *fiber.Ctx) error {
	url := googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	return c.Redirect(url, fiber.StatusTemporaryRedirect)
}

func (AuthController) HandleGoogleCallback(c *fiber.Ctx) error {
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(c.Context(), code)
	if err != nil {
		log.Error("Failed to exchange token:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to exchange token")
	}

	resp, err := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		log.Error("Failed to get user info:", err)
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user info")
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		log.Error("Failed to decode user info:", err)
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
	} else {
		params := db.CreateUserParams{
			Email: email,
			Token: token.RefreshToken,
		}
		user, err = connection_manager.Mgr.Queries.CreateUser(c.Context(), params)
		if err != nil {
			log.Error("Failed to create user:", err)
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	// Create the Claims
	claims := JWTClaims{
		ID: json.Number(fmt.Sprintf("%d", user.ID)),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token
	t, err := jwtToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		log.Error("Failed to sign token:", err)
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Return an HTML page that sends the token to the parent window
	html := fmt.Sprintf(`
		<html>
			<body>
				<script>
					window.opener.postMessage({ token: "%s" }, "http://localhost:5173");
				</script>
			</body>
		</html>
	`, t)

	return c.Type("html").Send([]byte(html))
}
