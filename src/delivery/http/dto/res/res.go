package res

import "github.com/turistikrota/service.upload/src/app/command"

type Response interface {
	ImageUploaded(res *command.UploadImageResult) *FileUploadedResponse
	PdfUploaded(res *command.UploadPdfResult) *FileUploadedResponse
	SvgUploaded(res *command.UploadSvgResult) *FileUploadedResponse
	MarkdownUploaded(res *command.UploadMarkdownResult) *FileUploadedResponse
	AvatarUploaded(res *command.UploadAvatarResult) *FileUploadedResponse
	BusinessAvatarUploaded(res *command.UploadBusinessAvatarResult) *FileUploadedResponse
	BusinessCoverUploaded(res *command.UploadBusinessCoverResult) *FileUploadedResponse
}

type response struct{}

func New() Response {
	return &response{}
}
