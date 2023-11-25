package command

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/mixarchitecture/i18np"
	"github.com/mixarchitecture/microp/decorator"
	"github.com/turistikrota/service.upload/src/domain/cdn"
)

type UploadBusinessAvatarCommand struct {
	NickName string
	Content  *multipart.FileHeader
}

type UploadBusinessAvatarResult struct {
	Url string
}

type UploadBusinessAvatarHandler decorator.CommandHandler[UploadBusinessAvatarCommand, *UploadBusinessAvatarResult]

type uploadBusinessAvatarHandler struct {
	repo    cdn.Repository
	factory cdn.Factory
}

type UploadBusinessAvatarHandlerConfig struct {
	Repo     cdn.Repository
	Factory  cdn.Factory
	CqrsBase decorator.Base
}

func NewUploadBusinessAvatarHandler(config UploadBusinessAvatarHandlerConfig) UploadBusinessAvatarHandler {
	return decorator.ApplyCommandDecorators[UploadBusinessAvatarCommand, *UploadBusinessAvatarResult](
		uploadBusinessAvatarHandler{
			repo:    config.Repo,
			factory: config.Factory,
		},
		config.CqrsBase,
	)
}

func (h uploadBusinessAvatarHandler) Handle(ctx context.Context, command UploadBusinessAvatarCommand) (*UploadBusinessAvatarResult, *i18np.Error) {
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
	return &UploadBusinessAvatarResult{Url: url}, nil
}
