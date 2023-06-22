package http

import (
	"time"

	"api.turistikrota.com/shared/auth/session"
	"api.turistikrota.com/shared/auth/token"
	"api.turistikrota.com/shared/server/http"
	"api.turistikrota.com/shared/server/http/auth"
	"api.turistikrota.com/shared/server/http/auth/current_account"
	"api.turistikrota.com/shared/server/http/auth/current_user"
	"api.turistikrota.com/shared/server/http/auth/device_uuid"
	"api.turistikrota.com/shared/server/http/auth/required_access"
	"api.turistikrota.com/shared/validator"
	"api.turistikrota.com/upload/src/app"
	"api.turistikrota.com/upload/src/config"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/fiber/v2/middleware/timeout"
	"github.com/mixarchitecture/i18np"
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
	router.Post("/image", h.isUploadAdminRole(Fields.Image), h.wrapWithTimeout(h.UploadImage))
	router.Post("/pdf", h.isUploadAdminRole(Fields.Pdf), h.wrapWithTimeout(h.UploadPdf))
	router.Post("/svg", h.isUploadAdminRole(Fields.Svg), h.wrapWithTimeout(h.UploadSvg))
	router.Post("/md", h.isUploadAdminRole(Fields.Markdown), h.wrapWithTimeout(h.UploadMarkdown))
	router.Post("/@:userName", h.currentAccount(), h.wrapWithTimeout(h.UploadAvatar))
	return router
}

func (h Server) deviceUUID() fiber.Handler {
	return device_uuid.New(device_uuid.Config{
		Domain: h.httpHeaders.Domain,
	})
}

func (h Server) currentAccount() fiber.Handler {
	return current_account.New(current_account.Config{
		FieldName: "userName",
		I18n: 	&h.i18n,
		RequiredKey: Messages.Error.RequiredAuth,
		ForbiddenKey: Messages.Error.Forbidden,
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

func (h Server) wrapWithTimeout(fn fiber.Handler) fiber.Handler {
	return timeout.NewWithContext(fn, 10*time.Second)
}
