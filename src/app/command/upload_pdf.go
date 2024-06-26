package command

import (
	"context"
	"mime/multipart"

	"github.com/mixarchitecture/i18np"
	"github.com/mixarchitecture/microp/decorator"
	"github.com/turistikrota/service.upload/src/domain/cdn"
)

type UploadPdfCommand struct {
	RandomName bool                  `json:"randomName"`
	Slugify    bool                  `json:"slugify"`
	FileName   string                `json:"fileName"`
	Dir        string                `json:"dir"`
	Content    *multipart.FileHeader `json:"content"`
	IsAdmin    bool                  `json:"-"`
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
	name := h.factory.GenerateName(command.FileName, command.RandomName, command.Slugify)
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
