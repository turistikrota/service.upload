package service

import (
	"api.turistikrota.com/upload/src/adapters/cloudflare"
	"api.turistikrota.com/upload/src/app"
	"api.turistikrota.com/upload/src/app/command"
	"api.turistikrota.com/upload/src/config"
	"api.turistikrota.com/upload/src/domain/cdn"
	"github.com/turistikrota/service.shared/decorator"
	"github.com/turistikrota/service.shared/validator"
)

type Config struct {
	App       config.App
	Validator *validator.Validator
}

func NewApplication(cnf Config) app.Application {
	cdnFactory := cdn.NewFactory()
	/*
	cdnRepo := bunny.New(bunny.Config{
		CdnHost:     cnf.App.CDN.Host,
		UploadHost:  cnf.App.CDN.UploadHost,
		StorageZone: cnf.App.CDN.StorageZone,
		ApiKey:      cnf.App.CDN.ApiKey,
	})
	*/
	cdnRepo := cloudflare.NewR2(cloudflare.Config{
		AccountId: cnf.App.R2.AccountId,
		AccessKey: cnf.App.R2.AccessKey,
		SecretKey: cnf.App.R2.SecretKey,
		Bucket:   cnf.App.R2.Bucket,
		PublicHost: cnf.App.R2.PublicHost,
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
