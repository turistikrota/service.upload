package cdn

import (
	"strings"

	"github.com/mixarchitecture/i18np"
)

type Errors interface {
	InternalError() *i18np.Error
	ContentRequired() *i18np.Error
	TypeNotAccepted([]string, string) *i18np.Error
	SizeTooBig(int64) *i18np.Error
	SizeTooSmall(int64) *i18np.Error
}

type cdnErrors struct{}

func newCdnErrors() Errors {
	return &cdnErrors{}
}

func (e *cdnErrors) InternalError() *i18np.Error {
	return i18np.NewError(I18nMessages.InternalError)
}

func (e *cdnErrors) ContentRequired() *i18np.Error {
	return i18np.NewError(I18nMessages.ContentRequired)
}

func (e *cdnErrors) TypeNotAccepted(accepts []string, fileType string) *i18np.Error {
	acc := strings.Join(accepts, ", ")
	return i18np.NewError(I18nMessages.TypeNotAccepted, i18np.P{
		"Accept": acc,
		"Type":   fileType,
	})
}

func (e *cdnErrors) SizeTooBig(size int64) *i18np.Error {
	return i18np.NewError(I18nMessages.SizeTooBig, i18np.P{
		"Size": size,
	})
}

func (e *cdnErrors) SizeTooSmall(size int64) *i18np.Error {
	return i18np.NewError(I18nMessages.SizeTooSmall, i18np.P{
		"Size": size,
	})
}
