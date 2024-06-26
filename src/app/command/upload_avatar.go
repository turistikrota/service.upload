package command

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/mixarchitecture/i18np"
	"github.com/mixarchitecture/microp/decorator"
	"github.com/turistikrota/service.upload/src/domain/cdn"
)

type UploadAvatarCommand struct {
	UserName string
	Content  *multipart.FileHeader
}

type UploadAvatarResult struct {
	Url string
}

type UploadAvatarHandler decorator.CommandHandler[UploadAvatarCommand, *UploadAvatarResult]

type uploadAvatarHandler struct {
	repo    cdn.Repository
	factory cdn.Factory
}

type UploadAvatarHandlerConfig struct {
	Repo     cdn.Repository
	Factory  cdn.Factory
	CqrsBase decorator.Base
}

func NewUploadAvatarHandler(config UploadAvatarHandlerConfig) UploadAvatarHandler {
	return decorator.ApplyCommandDecorators[UploadAvatarCommand, *UploadAvatarResult](
		uploadAvatarHandler{
			repo:    config.Repo,
			factory: config.Factory,
		},
		config.CqrsBase,
	)
}

func (h uploadAvatarHandler) Handle(ctx context.Context, command UploadAvatarCommand) (*UploadAvatarResult, *i18np.Error) {
	dir := "avatars"
	name := fmt.Sprintf("@%s", command.UserName)
	bytes, err := h.factory.NewImage(cdn.ValidateConfig{
		Content:     command.Content,
		Accept:      []string{"image/png", "image/jpeg", "image/jpg"},
		MaxSize:     1 * 1024 * 1024,
		MinSize:     1,
		Width:       1000,
		Height:      1000,
		MinifyLevel: cdn.MinifyLevelMedium,
	})
	if err != nil {
		return nil, err
	}
	fullName := name + ".png"
	url, success := h.repo.Upload(bytes, fullName, dir)
	if !success {
		return nil, h.factory.Errors.InternalError()
	}
	return &UploadAvatarResult{Url: url}, nil
}
