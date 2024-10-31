package controllers

import (
	"encoding/base64"
	"os"

	connection_manager "github.com/chenshuiluke/di-pocket-watcher/api/internal"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

type TransactionAnalysisController struct {
}

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

func (TransactionAnalysisController) GetTransactionEmails(c *fiber.Ctx) error {
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

	user, err := connection_manager.Mgr.Queries.GetUserToken(c.Context(), id)
	if err != nil {
		log.Error("Failed to get user:", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}
	updatedToken, err := googleOauthConfig.TokenSource(c.Context(), &oauth2.Token{RefreshToken: user.Token}).Token()
	client := googleOauthConfig.Client(c.Context(), updatedToken)
	gmailService, err := gmail.NewService(c.Context(), option.WithHTTPClient(client))
	if err != nil {
		log.Error("Failed to create gmail service", err)
	}
	response, err := gmailService.Users.Messages.List("me").Do()
	if err != nil {
		log.Error("Failed to retrieve the emails")
	}
	for _, record := range response.Messages {
		email, err := gmailService.Users.Messages.Get("me", record.Id).Format("full").Do()
		if err != nil || email == nil {
			log.Error("Failed to retrieve email:", err)
			continue
		}

		// Get the email body
		var body string
		if len(email.Payload.Parts) > 0 {
			// Multipart message
			for _, part := range email.Payload.Parts {
				if part.MimeType == "text/plain" {
					if data, err := base64.URLEncoding.DecodeString(part.Body.Data); err == nil {
						body = string(data)
						break
					}
				}
			}
		} else if email.Payload.Body.Data != "" {
			// Single part message
			if data, err := base64.URLEncoding.DecodeString(email.Payload.Body.Data); err == nil {
				body = string(data)
			}
		}

		log.Info("Email body:", body)
	}
	return nil
}
