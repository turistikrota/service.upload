package http

import (
	"fmt"

	"api.turistikrota.com/shared/server/http/auth/claim_guard"
	"api.turistikrota.com/shared/server/http/auth/current_user"
	"api.turistikrota.com/upload/src/config"
	"github.com/gofiber/fiber/v2"
)

func (h Server) isUploadAdminRole(field string) fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		u := current_user.Parse(ctx)
		role := fmt.Sprintf("%s.%s", config.Roles.Cdn.Upload, field)
		if u != nil && claim_guard.CheckClaim(u, role) {
			ctx.Locals(field, "true")
		}
		return ctx.Next()
	}
}
