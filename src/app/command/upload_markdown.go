package command

import (
	"context"
	"mime/multipart"

	"api.turistikrota.com/upload/src/domain/cdn"
	"github.com/mixarchitecture/i18np"
	"github.com/turistikrota/service.shared/decorator"
)

type UploadMarkdownCommand struct {
	RandomName bool
	FileName   string
	Dir        string
	Content    *multipart.FileHeader
	IsAdmin    bool
}

type UploadMarkdownResult struct {
	Url string
}

type UploadMarkdownHandler decorator.CommandHandler[UploadMarkdownCommand, *UploadMarkdownResult]

type uploadMarkdownHandler struct {
	repo    cdn.Repository
	factory cdn.Factory
}

type UploadMarkdownHandlerConfig struct {
	Repo     cdn.Repository
	Factory  cdn.Factory
	CqrsBase decorator.Base
}

func NewUploadMarkdownHandler(config UploadMarkdownHandlerConfig) UploadMarkdownHandler {
	return decorator.ApplyCommandDecorators[UploadMarkdownCommand, *UploadMarkdownResult](
		uploadMarkdownHandler{
			repo:    config.Repo,
			factory: config.Factory,
		},
		config.CqrsBase,
	)
}

func (h uploadMarkdownHandler) Handle(ctx context.Context, command UploadMarkdownCommand) (*UploadMarkdownResult, *i18np.Error) {
	dir := h.factory.GenerateDirName(command.Dir, command.IsAdmin, "md")
	name := h.factory.GenerateName(command.FileName, command.RandomName)
	bytes, err := h.factory.New(cdn.ValidateConfig{
		Content: command.Content,
		Accept:  []string{"text/markdown"},
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
	return &UploadMarkdownResult{Url: url}, nil
}
