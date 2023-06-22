package delivery

import (
	"api.turistikrota.com/upload/src/app"
	"api.turistikrota.com/upload/src/config"
	"api.turistikrota.com/upload/src/delivery/http"
	"github.com/gofiber/fiber/v2"
	"github.com/mixarchitecture/i18np"
	"github.com/turistikrota/service.shared/auth/session"
	"github.com/turistikrota/service.shared/auth/token"
	sharedHttp "github.com/turistikrota/service.shared/server/http"
	"github.com/turistikrota/service.shared/validator"
)

type Delivery interface {
	Load()
}

type delivery struct {
	app        app.Application
	config     config.App
	i18n       *i18np.I18n
	validator  *validator.Validator
	sessionSrv session.Service
	tknSrv     token.Service
}

type Config struct {
	App        app.Application
	Config     config.App
	I18n       *i18np.I18n
	Validator  *validator.Validator
	SessionSrv session.Service
	TokenSrv   token.Service
}

func New(config Config) Delivery {
	return &delivery{
		app:        config.App,
		config:     config.Config,
		i18n:       config.I18n,
		validator:  config.Validator,
		sessionSrv: config.SessionSrv,
		tknSrv:     config.TokenSrv,
	}
}

func (d *delivery) Load() {
	d.loadHTTP()
}

func (d *delivery) loadHTTP() *delivery {
	sharedHttp.RunServer(sharedHttp.Config{
		Host:  d.config.Server.Host,
		Port:  d.config.Server.Port,
		I18n:  d.i18n,
		Group: d.config.Server.Group,
		CreateHandler: func(router fiber.Router) fiber.Router {
			return http.New(http.Config{
				App:         d.app,
				I18n:        *d.i18n,
				Validator:   *d.validator,
				HttpHeaders: d.config.HttpHeaders,
				SessionSrv:  d.sessionSrv,
				TokenSrv:    d.tknSrv,
			}).Load(router)
		},
	})
	return d
}
