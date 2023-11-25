package res

import "github.com/turistikrota/service.upload/src/app/command"

type FileUploadedResponse struct {
	Url string `json:"url"`
}

func (r *response) ImageUploaded(res *command.UploadImageResult) *FileUploadedResponse {
	return &FileUploadedResponse{
		Url: res.Url,
	}
}

func (r *response) PdfUploaded(res *command.UploadPdfResult) *FileUploadedResponse {
	return &FileUploadedResponse{
		Url: res.Url,
	}
}

func (r *response) SvgUploaded(res *command.UploadSvgResult) *FileUploadedResponse {
	return &FileUploadedResponse{
		Url: res.Url,
	}
}

func (r *response) MarkdownUploaded(res *command.UploadMarkdownResult) *FileUploadedResponse {
	return &FileUploadedResponse{
		Url: res.Url,
	}
}

func (r *response) AvatarUploaded(res *command.UploadAvatarResult) *FileUploadedResponse {
	return &FileUploadedResponse{
		Url: res.Url,
	}
}

func (r *response) BusinessAvatarUploaded(res *command.UploadBusinessAvatarResult) *FileUploadedResponse {
	return &FileUploadedResponse{
		Url: res.Url,
	}
}

func (r *response) BusinessCoverUploaded(res *command.UploadBusinessCoverResult) *FileUploadedResponse {
	return &FileUploadedResponse{
		Url: res.Url,
	}
}
