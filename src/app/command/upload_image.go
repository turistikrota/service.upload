package command

import (
	"context"
	"mime/multipart"

	"api.turistikrota.com/upload/src/domain/cdn"
	"github.com/mixarchitecture/i18np"
	"github.com/turistikrota/service.shared/decorator"
)

type UploadImageCommand struct {
	RandomName bool
	FileName   string
	Dir        string
	Content    *multipart.FileHeader
	IsAdmin    bool
}

type UploadImageResult struct {
	Url string
}

type UploadImageHandler decorator.CommandHandler[UploadImageCommand, *UploadImageResult]

type uploadImageHandler struct {
	repo    cdn.Repository
	factory cdn.Factory
}

type UploadImageHandlerConfig struct {
	Repo     cdn.Repository
	Factory  cdn.Factory
	CqrsBase decorator.Base
}

func NewUploadImageHandler(config UploadImageHandlerConfig) UploadImageHandler {
	return decorator.ApplyCommandDecorators[UploadImageCommand, *UploadImageResult](
		uploadImageHandler{
			repo:    config.Repo,
			factory: config.Factory,
		},
		config.CqrsBase,
	)
}

func (h uploadImageHandler) Handle(ctx context.Context, command UploadImageCommand) (*UploadImageResult, *i18np.Error) {
	dir := h.factory.GenerateDirName(command.Dir, command.IsAdmin, "img")
	name := h.factory.GenerateName(command.FileName, command.RandomName)
	bytes, err := h.factory.NewImage(cdn.ValidateConfig{
		Content: command.Content,
		Accept:  []string{"image/jpeg", "image/jpg", "image/png", "image/gif"},
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
	return &UploadImageResult{Url: url}, nil
}
