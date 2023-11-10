package command

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/mixarchitecture/i18np"
	"github.com/mixarchitecture/microp/decorator"
	"github.com/turistikrota/service.upload/src/domain/cdn"
)

type UploadOwnerAvatarCommand struct {
	NickName string
	Content  *multipart.FileHeader
}

type UploadOwnerAvatarResult struct {
	Url string
}

type UploadOwnerAvatarHandler decorator.CommandHandler[UploadOwnerAvatarCommand, *UploadOwnerAvatarResult]

type uploadOwnerAvatarHandler struct {
	repo    cdn.Repository
	factory cdn.Factory
}

type UploadOwnerAvatarHandlerConfig struct {
	Repo     cdn.Repository
	Factory  cdn.Factory
	CqrsBase decorator.Base
}

func NewUploadOwnerAvatarHandler(config UploadOwnerAvatarHandlerConfig) UploadOwnerAvatarHandler {
	return decorator.ApplyCommandDecorators[UploadOwnerAvatarCommand, *UploadOwnerAvatarResult](
		uploadOwnerAvatarHandler{
			repo:    config.Repo,
			factory: config.Factory,
		},
		config.CqrsBase,
	)
}

func (h uploadOwnerAvatarHandler) Handle(ctx context.Context, command UploadOwnerAvatarCommand) (*UploadOwnerAvatarResult, *i18np.Error) {
	dir := "avatars"
	name := fmt.Sprintf("~%s", command.NickName)
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
	return &UploadOwnerAvatarResult{Url: url}, nil
}
