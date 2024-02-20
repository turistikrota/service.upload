package main

import (
	"github.com/mixarchitecture/microp/events/nats"
	"github.com/mixarchitecture/microp/validator"
	"github.com/turistikrota/service.shared/auth/session"
	"github.com/turistikrota/service.shared/auth/token"
	"github.com/turistikrota/service.shared/db/redis"

	"github.com/mixarchitecture/i18np"
	"github.com/mixarchitecture/microp/env"
	"github.com/mixarchitecture/microp/logs"
	"github.com/turistikrota/service.upload/src/config"
	"github.com/turistikrota/service.upload/src/delivery"
	"github.com/turistikrota/service.upload/src/service"
)

func main() {
	logs.Init()
	cnf := config.App{}
	env.Load(&cnf)
	i18n := i18np.New(cnf.I18n.Fallback)
	i18n.Load(cnf.I18n.Dir, cnf.I18n.Locales...)
	eventEngine := nats.New(nats.Config{
		Url:     cnf.Nats.Url,
		Streams: cnf.Nats.Streams,
	})
	valid := validator.New(i18n)
	valid.ConnectCustom()
	valid.RegisterTagName()
	app := service.NewApplication(service.Config{
		App:       cnf,
		Validator: valid,
	})
	r := redis.New(&redis.Config{
		Host:     cnf.Redis.Host,
		Port:     cnf.Redis.Port,
		Password: cnf.Redis.Pw,
		DB:       cnf.Redis.Db,
	})
	tknSrv := token.New(token.Config{
		Expiration:     cnf.TokenSrv.Expiration,
		PublicKeyFile:  cnf.RSA.PublicKeyFile,
		PrivateKeyFile: cnf.RSA.PrivateKeyFile,
	})
	session := session.NewSessionApp(session.Config{
		Redis:       r,
		EventEngine: eventEngine,
		TokenSrv:    tknSrv,
		Topic:       cnf.Session.Topic,
		Project:     cnf.TokenSrv.Project,
	})
	delivery := delivery.New(delivery.Config{
		App:         app,
		Config:      cnf,
		I18n:        i18n,
		Validator:   valid,
		SessionSrv:  session.Service,
		TokenSrv:    tknSrv,
		EventEngine: eventEngine,
	})
	delivery.Load()
}
