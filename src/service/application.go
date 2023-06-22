package service

import (
	"api.turistikrota.com/shared/decorator"
	"api.turistikrota.com/shared/validator"
	"api.turistikrota.com/upload/src/adapters/bunny"
	"api.turistikrota.com/upload/src/app"
	"api.turistikrota.com/upload/src/app/command"
	"api.turistikrota.com/upload/src/config"
	"api.turistikrota.com/upload/src/domain/cdn"
)

type Config struct {
	App       config.App
	Validator *validator.Validator
}

func NewApplication(cnf Config) app.Application {
	cdnFactory := cdn.NewFactory()
	cdnRepo := bunny.New(bunny.Config{
		CdnHost:     cnf.App.CDN.Host,
		UploadHost:  cnf.App.CDN.UploadHost,
		StorageZone: cnf.App.CDN.StorageZone,
		ApiKey:      cnf.App.CDN.ApiKey,
	})

	base := decorator.NewBase()

	return app.Application{
		Commands: app.Commands{
			UploadImage: command.NewUploadImageHandler(command.UploadImageHandlerConfig{
				Repo:     cdnRepo,
				Factory:  cdnFactory,
				CqrsBase: base,
			}),
			UploadAvatar: command.NewUploadAvatarHandler(command.UploadAvatarHandlerConfig{
				Repo:     cdnRepo,
				Factory:  cdnFactory,
				CqrsBase: base,
			}),
			UploadMarkdown: command.NewUploadMarkdownHandler(command.UploadMarkdownHandlerConfig{
				Repo:     cdnRepo,
				Factory:  cdnFactory,
				CqrsBase: base,
			}),
			UploadPdf: command.NewUploadPdfHandler(command.UploadPdfHandlerConfig{
				Repo:     cdnRepo,
				Factory:  cdnFactory,
				CqrsBase: base,
			}),
			UploadSvg: command.NewUploadSvgHandler(command.UploadSvgHandlerConfig{
				Repo:     cdnRepo,
				Factory:  cdnFactory,
				CqrsBase: base,
			}),
		},
		Queries: app.Queries{},
	}
}
