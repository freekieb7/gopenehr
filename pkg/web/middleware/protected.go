package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type ValidateTokenFunc func(ctx context.Context, token string) (map[string]any, error)

func JWTProtected(scopes []string, validate ValidateTokenFunc) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")
		if tokenString == "" {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		if tokenString == "" {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		claims, err := validate(c.Context(), tokenString)
		if err != nil {
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		if len(scopes) > 0 {
			tokenScopesRaw, ok := claims["scope"]
			if !ok {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			tokenScopesStr, ok := tokenScopesRaw.(string)
			if !ok {
				c.Status(fiber.StatusForbidden)
				return nil
			}

			tokenScopes := strings.Split(tokenScopesStr, " ")
			scopeMap := make(map[string]bool)
			for _, s := range tokenScopes {
				scopeMap[s] = true
			}

			for _, requiredScope := range scopes {
				if !scopeMap[string(requiredScope)] {
					c.Status(fiber.StatusForbidden)
					return nil
				}
			}
		}

		return c.Next()
	}
}
