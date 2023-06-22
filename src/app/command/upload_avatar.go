package command

import (
	"context"
	"fmt"
	"mime/multipart"

	"api.turistikrota.com/upload/src/domain/cdn"
	"github.com/mixarchitecture/i18np"
	"github.com/turistikrota/service.shared/decorator"
)

type UploadAvatarCommand struct {
	UserName string
	UserCode string
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
	dir := "avatar"
	name := fmt.Sprintf("%s-%s", command.UserName, command.UserCode)
	bytes, err := h.factory.New(cdn.ValidateConfig{
		Content: command.Content,
		Accept:  []string{"image/png"},
		MaxSize: 3 * 1024 * 1024,
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
	return &UploadAvatarResult{Url: url}, nil
}
