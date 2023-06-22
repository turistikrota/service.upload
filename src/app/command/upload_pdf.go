package command

import (
	"context"
	"mime/multipart"

	"api.turistikrota.com/upload/src/domain/cdn"
	"github.com/mixarchitecture/i18np"
	"github.com/turistikrota/service.shared/decorator"
)

type UploadPdfCommand struct {
	RandomName bool
	FileName   string
	Dir        string
	Content    *multipart.FileHeader
	IsAdmin    bool
}

type UploadPdfResult struct {
	Url string
}

type UploadPdfHandler decorator.CommandHandler[UploadPdfCommand, *UploadPdfResult]

type uploadPdfHandler struct {
	repo    cdn.Repository
	factory cdn.Factory
}

type UploadPdfHandlerConfig struct {
	Repo     cdn.Repository
	Factory  cdn.Factory
	CqrsBase decorator.Base
}

func NewUploadPdfHandler(config UploadPdfHandlerConfig) UploadPdfHandler {
	return decorator.ApplyCommandDecorators[UploadPdfCommand, *UploadPdfResult](
		uploadPdfHandler{
			repo:    config.Repo,
			factory: config.Factory,
		},
		config.CqrsBase,
	)
}

func (h uploadPdfHandler) Handle(ctx context.Context, command UploadPdfCommand) (*UploadPdfResult, *i18np.Error) {
	dir := h.factory.GenerateDirName(command.Dir, command.IsAdmin, "pdf")
	name := h.factory.GenerateName(command.FileName, command.RandomName)
	bytes, err := h.factory.New(cdn.ValidateConfig{
		Content: command.Content,
		Accept:  []string{"application/pdf"},
		MaxSize: 5 * 1024 * 1024,
		MinSize: 1,
	})
	if err != nil {
		return nil, err
	}
	fullName := name + "." + h.factory.GetExtension(command.Content)
	url, success := h.repo.Upload(bytes, fullName, dir)
	if !success {
		return nil, h.factory.Errors.InternalError()
	}
	return &UploadPdfResult{Url: url}, nil
}
