package command

import (
	"context"
	"mime/multipart"

	"github.com/mixarchitecture/i18np"
	"github.com/mixarchitecture/microp/decorator"
	"github.com/turistikrota/service.upload/src/domain/cdn"
)

type UploadImageCommand struct {
	RandomName  bool                  `json:"randomName"`
	Slugify     bool                  `json:"slugify"`
	FileName    string                `json:"fileName"`
	Dir         string                `json:"dir"`
	Content     *multipart.FileHeader `json:"content"`
	IsAdmin     bool                  `json:"-"`
	MinifyLevel cdn.MinifyLevel       `json:"minifyLevel"`
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
	name := h.factory.GenerateName(command.FileName, command.RandomName, command.Slugify)
	bytes, err := h.factory.NewImage(cdn.ValidateConfig{
		Content:     command.Content,
		Accept:      []string{"image/jpeg", "image/jpg", "image/png", "image/gif", "image/webp"},
		MaxSize:     5 * 1024 * 1024,
		MinSize:     1,
		Width:       1000,
		Height:      1000,
		MinifyLevel: command.MinifyLevel,
	})
	if err != nil {
		return nil, err
	}
	fullName := name + ".webp"
	url, success := h.repo.Upload(bytes, fullName, dir)
	if !success {
		return nil, h.factory.Errors.InternalError()
	}
	return &UploadImageResult{Url: url}, nil
}
