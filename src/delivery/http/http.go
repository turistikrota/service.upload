package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/mixarchitecture/i18np"
	"github.com/mixarchitecture/microp/server/http"
	"github.com/mixarchitecture/microp/validator"
	"github.com/turistikrota/service.shared/auth/session"
	"github.com/turistikrota/service.shared/auth/token"
	"github.com/turistikrota/service.shared/server/http/auth"
	"github.com/turistikrota/service.shared/server/http/auth/current_account"
	"github.com/turistikrota/service.shared/server/http/auth/current_owner"
	"github.com/turistikrota/service.shared/server/http/auth/current_user"
	"github.com/turistikrota/service.shared/server/http/auth/device_uuid"
	"github.com/turistikrota/service.shared/server/http/auth/required_access"
	"github.com/turistikrota/service.upload/src/app"
	"github.com/turistikrota/service.upload/src/config"
)

type Server struct {
	app         app.Application
	i18n        i18np.I18n
	validator   validator.Validator
	httpHeaders config.HttpHeaders
	tknSrv      token.Service
	sessionSrv  session.Service
}

type Config struct {
	App         app.Application
	I18n        i18np.I18n
	Validator   validator.Validator
	HttpHeaders config.HttpHeaders
	TokenSrv    token.Service
	SessionSrv  session.Service
}

func New(config Config) Server {
	return Server{
		app:         config.App,
		i18n:        config.I18n,
		validator:   config.Validator,
		httpHeaders: config.HttpHeaders,
		tknSrv:      config.TokenSrv,
		sessionSrv:  config.SessionSrv,
	}
}

func (h Server) Load(router fiber.Router) fiber.Router {
	router.Use(h.cors(), h.deviceUUID(), h.rateLimit(), h.currentUserAccess(), h.requiredAccess())
	router.Post("/image", h.isUploadAdminRole(Fields.Image), h.wrapWithTimeout(h.UploadImage, 3*time.Second))
	router.Post("/pdf", h.isUploadAdminRole(Fields.Pdf), h.wrapWithTimeout(h.UploadPdf))
	router.Post("/svg", h.isUploadAdminRole(Fields.Svg), h.wrapWithTimeout(h.UploadSvg))
	router.Post("/md", h.isUploadAdminRole(Fields.Markdown), h.wrapWithTimeout(h.UploadMarkdown))
	router.Post("/@:userName", h.currentAccount(), h.wrapWithTimeout(h.UploadAvatar, 3*time.Second))
	router.Post("/owner/avatar", h.currentOwner(config.Roles.Owner.Super, config.Roles.Owner.UploadAvatar), h.wrapWithTimeout(h.UploadOwnerAvatar, 3*time.Second))
	router.Post("/owner/cover", h.currentOwner(config.Roles.Owner.Super, config.Roles.Owner.UploadCover), h.wrapWithTimeout(h.UploadOwnerCover, 3*time.Second))
	return router
}

func (h Server) deviceUUID() fiber.Handler {
	return device_uuid.New(device_uuid.Config{
		Domain: h.httpHeaders.Domain,
	})
}

func (h Server) currentAccount() fiber.Handler {
	return current_account.New(current_account.Config{
		FieldName:    "userName",
		I18n:         &h.i18n,
		RequiredKey:  Messages.Error.RequiredAuth,
		ForbiddenKey: Messages.Error.Forbidden,
	})
}

func (h Server) currentOwner(roles ...string) fiber.Handler {
	return current_owner.New(current_owner.Config{
		I18n:         &h.i18n,
		Roles:        roles,
		RequiredKey:  Messages.Error.RequiredOwnerSelect,
		ForbiddenKey: Messages.Error.ForbiddenOwnerSelect,
	})
}

func (h Server) currentUserAccess() fiber.Handler {
	return current_user.New(current_user.Config{
		TokenSrv:   h.tknSrv,
		SessionSrv: h.sessionSrv,
		I18n:       &h.i18n,
		MsgKey:     Messages.Error.CurrentUserAccess,
		HeaderKey:  http.Headers.Authorization,
		CookieKey:  auth.Cookies.AccessToken,
		UseCookie:  true,
		UseBearer:  true,
		IsRefresh:  false,
		IsAccess:   true,
	})
}

func (h Server) requiredAccess() fiber.Handler {
	return required_access.New(required_access.Config{
		I18n:   h.i18n,
		MsgKey: Messages.Error.RequiredAuth,
	})
}

func (h Server) cors() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     h.httpHeaders.AllowedOrigins,
		AllowMethods:     h.httpHeaders.AllowedMethods,
		AllowHeaders:     h.httpHeaders.AllowedHeaders,
		AllowCredentials: h.httpHeaders.AllowCredentials,
	})
}

func (h Server) rateLimit() fiber.Handler {
	return limiter.New(limiter.Config{
		Max:        15,
		Expiration: 3 * time.Minute,
	})
}

func (h Server) wrapWithTimeout(fn fiber.Handler, secs ...time.Duration) fiber.Handler {
	sec := 1 * time.Second
	if len(secs) > 0 {
		sec = secs[0]
	}
	return timeout.NewWithContext(fn, sec*time.Second)
}
