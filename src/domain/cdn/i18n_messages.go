package cdn

type messages struct {
	InternalError   string
	ContentRequired string
	TypeNotAccepted string
	SizeTooBig      string
	SizeTooSmall    string
}

var I18nMessages = messages{
	InternalError:   "cdn_internal_error",
	ContentRequired: "cdn_content_required",
	TypeNotAccepted: "cdn_type_not_accepted",
	SizeTooBig:      "cdn_size_too_big",
	SizeTooSmall:    "cdn_size_too_small",
}
