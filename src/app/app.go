package app

import (
	"github.com/turistikrota/service.upload/src/app/command"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	UploadImage          command.UploadImageHandler
	UploadMarkdown       command.UploadMarkdownHandler
	UploadPdf            command.UploadPdfHandler
	UploadSvg            command.UploadSvgHandler
	UploadAvatar         command.UploadAvatarHandler
	UploadBusinessAvatar command.UploadBusinessAvatarHandler
	UploadBusinessCover  command.UploadBusinessCoverHandler
}

type Queries struct{}
