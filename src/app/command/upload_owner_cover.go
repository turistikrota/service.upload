package command

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/mixarchitecture/i18np"
	"github.com/mixarchitecture/microp/decorator"
	"github.com/turistikrota/service.upload/src/domain/cdn"
)

type UploadBusinessCoverCommand struct {
	NickName string
	Content  *multipart.FileHeader
}

type UploadBusinessCoverResult struct {
	Url string
}

type UploadBusinessCoverHandler decorator.CommandHandler[UploadBusinessCoverCommand, *UploadBusinessCoverResult]

type uploadBusinessCoverHandler struct {
	repo    cdn.Repository
	factory cdn.Factory
}

type UploadBusinessCoverHandlerConfig struct {
	Repo     cdn.Repository
	Factory  cdn.Factory
	CqrsBase decorator.Base
}

func NewUploadBusinessCoverHandler(config UploadBusinessCoverHandlerConfig) UploadBusinessCoverHandler {
	return decorator.ApplyCommandDecorators[UploadBusinessCoverCommand, *UploadBusinessCoverResult](
		uploadBusinessCoverHandler{
			repo:    config.Repo,
			factory: config.Factory,
		},
		config.CqrsBase,
	)
}

func (h uploadBusinessCoverHandler) Handle(ctx context.Context, command UploadBusinessCoverCommand) (*UploadBusinessCoverResult, *i18np.Error) {
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
	return &UploadBusinessCoverResult{Url: url}, nil
}
