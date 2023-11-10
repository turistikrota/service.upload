package command

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/mixarchitecture/i18np"
	"github.com/mixarchitecture/microp/decorator"
	"github.com/turistikrota/service.upload/src/domain/cdn"
)

type UploadOwnerCoverCommand struct {
	NickName string
	Content  *multipart.FileHeader
}

type UploadOwnerCoverResult struct {
	Url string
}

type UploadOwnerCoverHandler decorator.CommandHandler[UploadOwnerCoverCommand, *UploadOwnerCoverResult]

type uploadOwnerCoverHandler struct {
	repo    cdn.Repository
	factory cdn.Factory
}

type UploadOwnerCoverHandlerConfig struct {
	Repo     cdn.Repository
	Factory  cdn.Factory
	CqrsBase decorator.Base
}

func NewUploadOwnerCoverHandler(config UploadOwnerCoverHandlerConfig) UploadOwnerCoverHandler {
	return decorator.ApplyCommandDecorators[UploadOwnerCoverCommand, *UploadOwnerCoverResult](
		uploadOwnerCoverHandler{
			repo:    config.Repo,
			factory: config.Factory,
		},
		config.CqrsBase,
	)
}

func (h uploadOwnerCoverHandler) Handle(ctx context.Context, command UploadOwnerCoverCommand) (*UploadOwnerCoverResult, *i18np.Error) {
	dir := "covers"
	name := fmt.Sprintf("~%s", command.NickName)
	bytes, err := h.factory.NewImage(cdn.ValidateConfig{
		Content:     command.Content,
		Accept:      []string{"image/png", "image/jpeg", "image/jpg"},
		MaxSize:     1 * 1024 * 1024,
		MinSize:     1,
		Width:       1500,
		Height:      500,
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
	return &UploadOwnerCoverResult{Url: url}, nil
}
