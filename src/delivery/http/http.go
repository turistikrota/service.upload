package http

import (
	"strings"
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
	"github.com/turistikrota/service.shared/server/http/auth/current_business"
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
	router.Use(h.cors(), h.deviceUUID(), h.currentUserAccess(), h.requiredAccess())
	router.Post("/image", h.rateLimit(75), h.isUploadAdminRole(Fields.Image), h.wrapWithTimeout(h.UploadImage, 3*time.Second))
	router.Post("/pdf", h.rateLimit(), h.isUploadAdminRole(Fields.Pdf), h.wrapWithTimeout(h.UploadPdf))
	router.Post("/svg", h.rateLimit(), h.isUploadAdminRole(Fields.Svg), h.wrapWithTimeout(h.UploadSvg))
	router.Post("/md", h.rateLimit(), h.isUploadAdminRole(Fields.Markdown), h.wrapWithTimeout(h.UploadMarkdown))
	router.Post("/@:userName", h.rateLimit(), h.currentAccount(), h.wrapWithTimeout(h.UploadAvatar, 3*time.Second))
	router.Post("/business/avatar", h.rateLimit(), h.currentAccount(), h.currentBusiness(config.Roles.Business.Super, config.Roles.Business.UploadAvatar), h.wrapWithTimeout(h.UploadBusinessAvatar, 3*time.Second))
	router.Post("/business/cover", h.rateLimit(), h.currentAccount(), h.currentBusiness(config.Roles.Business.Super, config.Roles.Business.UploadCover), h.wrapWithTimeout(h.UploadBusinessCover, 3*time.Second))
	router.Post("/listing", h.rateLimit(), h.currentAccount(), h.currentBusiness(config.Roles.Business.Super, config.Roles.Business.ListingCreate, config.Roles.Business.ListingUpdate), h.wrapWithTimeout(h.UploadListingImage, 3*time.Second))
	return router
}

func (h Server) deviceUUID() fiber.Handler {
	return device_uuid.New(device_uuid.Config{
		Domain: h.httpHeaders.Domain,
	})
}

func (h Server) currentAccount() fiber.Handler {
	return current_account.New(current_account.Config{
		I18n:         &h.i18n,
		RequiredKey:  Messages.Error.RequiredAuth,
		ForbiddenKey: Messages.Error.Forbidden,
	})
}

func (h Server) currentBusiness(roles ...string) fiber.Handler {
	return current_business.New(current_business.Config{
		I18n:         &h.i18n,
		Roles:        roles,
		RequiredKey:  Messages.Error.RequiredBusinessSelect,
		ForbiddenKey: Messages.Error.ForbiddenBusinessSelect,
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
		AllowMethods:     h.httpHeaders.AllowedMethods,
		AllowHeaders:     h.httpHeaders.AllowedHeaders,
		AllowCredentials: h.httpHeaders.AllowCredentials,
		AllowOriginsFunc: func(origin string) bool {
			origins := strings.Split(h.httpHeaders.AllowedOrigins, ",")
			for _, o := range origins {
				if strings.Contains(origin, o) {
					return true
				}
			}
			return false
		},
	})
}

func (h Server) rateLimit(limit ...int) fiber.Handler {
	max := 15
	if len(limit) > 0 {
		max = limit[0]
	}
	return limiter.New(limiter.Config{
		Max:        max,
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
