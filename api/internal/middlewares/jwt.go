package middlewares

import (
	"os"

	"github.com/chenshuiluke/di-pocket-watcher/api/internal/controllers"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func Jwt() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{Key: []byte(os.Getenv("JWT_SECRET"))},
		Claims:     &controllers.JWTClaims{},
	})
}
