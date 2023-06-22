package app

import (
	"api.turistikrota.com/upload/src/app/command"
)

type Application struct {
	Commands Commands
	Queries  Queries
}

type Commands struct {
	UploadImage    command.UploadImageHandler
	UploadMarkdown command.UploadMarkdownHandler
	UploadPdf      command.UploadPdfHandler
	UploadSvg      command.UploadSvgHandler
	UploadAvatar   command.UploadAvatarHandler
}

type Queries struct{}
