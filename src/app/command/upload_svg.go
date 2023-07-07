package command

import (
	"context"
	"mime/multipart"

	"github.com/mixarchitecture/i18np"
	"github.com/mixarchitecture/microp/decorator"
	"github.com/turistikrota/service.upload/src/domain/cdn"
)

type UploadSvgCommand struct {
	RandomName bool
	FileName   string
	Dir        string
	Content    *multipart.FileHeader
	IsAdmin    bool
}

type UploadSvgResult struct {
	Url string
}

type UploadSvgHandler decorator.CommandHandler[UploadSvgCommand, *UploadSvgResult]

type uploadSvgHandler struct {
	repo    cdn.Repository
	factory cdn.Factory
}

type UploadSvgHandlerConfig struct {
	Repo     cdn.Repository
	Factory  cdn.Factory
	CqrsBase decorator.Base
}

func NewUploadSvgHandler(config UploadSvgHandlerConfig) UploadSvgHandler {
	return decorator.ApplyCommandDecorators[UploadSvgCommand, *UploadSvgResult](
		uploadSvgHandler{
			repo:    config.Repo,
			factory: config.Factory,
		},
		config.CqrsBase,
	)
}

func (h uploadSvgHandler) Handle(ctx context.Context, command UploadSvgCommand) (*UploadSvgResult, *i18np.Error) {
	dir := h.factory.GenerateDirName(command.Dir, command.IsAdmin, "svg")
	name := h.factory.GenerateName(command.FileName, command.RandomName)
	bytes, err := h.factory.New(cdn.ValidateConfig{
		Content: command.Content,
		Accept:  []string{"image/svg+xml"},
		MaxSize: 5 * 1024 * 1024,
		MinSize: 1,
		Width:   1000,
		Height:  1000,
	})
	if err != nil {
		return nil, err
	}
	fullName := name + "." + h.factory.GetExtension(command.Content)
	url, success := h.repo.Upload(bytes, fullName, dir)
	if !success {
		return nil, h.factory.Errors.InternalError()
	}
	return &UploadSvgResult{Url: url}, nil
}
