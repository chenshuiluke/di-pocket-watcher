package controllers

import (
	"encoding/base64"
	"os"
	"strings"

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

var bankTransactionNotificationSenderEmailAddresses = []string{"no-reply-ncbcardalerts@jncb.com"}

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
		// Get only metadata first
		email, err := gmailService.Users.Messages.Get("me", record.Id).Format("metadata").Fields("payload/headers").Do()
		if err != nil || email == nil {
			log.Error("Failed to retrieve email metadata:", err)
			continue
		}

		// Find the "From" header
		var fromAddress string
		for _, header := range email.Payload.Headers {
			if header.Name == "From" {
				fromAddress = header.Value
				break
			}
		}

		// Check if the sender is in the allowed list
		senderFound := false
		for _, allowedSender := range bankTransactionNotificationSenderEmailAddresses {
			if strings.Contains(fromAddress, allowedSender) {
				senderFound = true
				break
			}
		}

		if !senderFound {
			continue
		}

		// Now get the full email only if it's from a bank sender
		fullEmail, err := gmailService.Users.Messages.Get("me", record.Id).Format("full").Do()
		if err != nil || fullEmail == nil {
			log.Error("Failed to retrieve full email:", err)
			continue
		}

		// Get the email body
		var body string
		if len(fullEmail.Payload.Parts) > 0 {
			// Multipart message
			for _, part := range fullEmail.Payload.Parts {
				if part.MimeType == "text/plain" {
					if data, err := base64.URLEncoding.DecodeString(part.Body.Data); err == nil {
						body = string(data)
						break
					}
				}
			}
		} else if fullEmail.Payload.Body.Data != "" {
			// Single part message
			if data, err := base64.URLEncoding.DecodeString(fullEmail.Payload.Body.Data); err == nil {
				body = string(data)
			}
		}

		log.Info("Email body:", body)
	}
	return nil
}
