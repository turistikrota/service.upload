package http

type successMessages struct {
	Create           string
	Get              string
	ImageUploaded    string
	PdfUploaded      string
	MarkdownUploaded string
	AvatarUploaded   string
}

type errorMessages struct {
	Create             string
	Get                string
	OnlyAcceptFormData string
	ImageNotFound      string
	PdfNotFound        string
	MarkdownNotFound   string
	RequiredAuth       string
	AdminRoute         string
	CurrentUserAccess  string
	UserNameRequired   string
	AvatarNotFound     string
	Forbidden          string
}

type messages struct {
	Success successMessages
	Error   errorMessages
}

var Messages = messages{
	Success: successMessages{
		Create:           "http_success_create",
		Get:              "http_success_get",
		ImageUploaded:    "http_success_image_uploaded",
		PdfUploaded:      "http_success_pdf_uploaded",
		MarkdownUploaded: "http_success_markdown_uploaded",
		AvatarUploaded:   "http_success_avatar_uploaded",
	},
	Error: errorMessages{
		Create:             "http_error_create",
		Get:                "http_error_get",
		OnlyAcceptFormData: "http_error_only_accept_form_data",
		ImageNotFound:      "http_error_image_not_found",
		PdfNotFound:        "http_error_pdf_not_found",
		MarkdownNotFound:   "http_error_markdown_not_found",
		RequiredAuth:       "http_required_auth",
		AdminRoute:         "http_admin_route",
		CurrentUserAccess:  "http_current_user_access",
		UserNameRequired:   "http_user_name_required",
		AvatarNotFound:     "http_avatar_not_found",
		Forbidden:          "http_forbidden",
	},
}
