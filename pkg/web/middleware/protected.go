package middleware

import (
	"context"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ValidateTokenFunc func(ctx context.Context, token string) (map[string]any, error)

func JWTProtected(scopes []string, validate ValidateTokenFunc) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		if validate == nil {
			return c.Next()
		}

		auditCtx := AuditFrom(c)

		tokenString := c.Get("Authorization")
		if tokenString == "" {
			auditCtx.Fail("unauthorized", "missing authorization header")
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		if tokenString == "" {
			auditCtx.Fail("unauthorized", "invalid authorization header")
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		claims, err := validate(c.Context(), tokenString)
		if err != nil {
			auditCtx.Fail("unauthorized", "invalid token")
			c.Status(fiber.StatusUnauthorized)
			return nil
		}

		if claims["tenant_id"] != nil {
			tenantIDStr, ok := claims["tenant_id"].(string)
			if !ok {
				auditCtx.Fail("unauthorized", "invalid tenant_id in token")
				c.Status(fiber.StatusUnauthorized)
				return nil
			}

			tenantID, err := uuid.Parse(tenantIDStr)
			if err != nil {
				auditCtx.Fail("unauthorized", "invalid tenant_id in token")
				c.Status(fiber.StatusUnauthorized)
				return nil
			}

			c.Locals("tenant_id", tenantID)
		}

		if len(scopes) > 0 {
			tokenScopesRaw, ok := claims["scope"]
			if !ok {
				auditCtx.Fail("forbidden", "missing scope in token")
				c.Status(fiber.StatusForbidden)
				return nil
			}

			tokenScopesStr, ok := tokenScopesRaw.(string)
			if !ok {
				auditCtx.Fail("forbidden", "invalid scope in token")
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
					auditCtx.Fail("forbidden", "insufficient scope in token")
					c.Status(fiber.StatusForbidden)
					return nil
				}
			}
		}

		return c.Next()
	}
}
