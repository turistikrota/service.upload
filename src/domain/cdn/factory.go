package cdn

import (
	"mime/multipart"

	"github.com/google/uuid"
	"github.com/mixarchitecture/i18np"
)

type Factory struct {
	Errors Errors
}

func NewFactory() Factory {
	return Factory{
		Errors: newCdnErrors(),
	}
}

func (f Factory) GenerateName(name string, random bool) string {
	if random || name == "" {
		return f.RandomName()
	}
	return name
}

func (f Factory) GenerateDirName(dir string, isAdmin bool, fb string) string {
	if isAdmin && dir != "" {
		return dir
	}
	return fb
}

func (f Factory) RandomName() string {
	return uuid.New().String()
}

type ValidateConfig struct {
	Content   *multipart.FileHeader
	Accept    []string
	MaxSize   int64
	MinSize   int64
	Width     int
	MinWidth  int
	MaxWidth  int
	Height    int
	MinHeight int
	MaxHeight int
}

func (f Factory) New(cnf ValidateConfig) ([]byte, *i18np.Error) {
	err := f.Validate(cnf)
	if err != nil {
		return nil, err
	}
	file, error := cnf.Content.Open()
	if error != nil {
		return nil, f.Errors.InternalError()
	}
	defer file.Close()
	bytes := make([]byte, cnf.Content.Size)
	file.Read(bytes)
	return bytes, nil
}

func (f Factory) Validate(cnf ValidateConfig) *i18np.Error {
	validators := []func(ValidateConfig) *i18np.Error{
		f.validateContent,
		f.validateAccept,
		f.validateSize,
		f.validateWidth,
		f.validateHeight,
	}
	for _, validator := range validators {
		err := validator(cnf)
		if err != nil {
			return err
		}
	}
	return nil
}

func (f Factory) validateContent(cnf ValidateConfig) *i18np.Error {
	if cnf.Content == nil {
		return f.Errors.ContentRequired()
	}
	return nil
}

func (f Factory) validateAccept(cnf ValidateConfig) *i18np.Error {
	if len(cnf.Accept) == 0 {
		return nil
	}
	for _, accept := range cnf.Accept {
		if accept == cnf.Content.Header.Get("Content-Type") {
			return nil
		}
	}
	return f.Errors.TypeNotAccepted(cnf.Accept)
}

func (f Factory) validateSize(cnf ValidateConfig) *i18np.Error {
	if cnf.MaxSize == 0 {
		return nil
	}
	if cnf.Content.Size > cnf.MaxSize {
		return f.Errors.SizeTooBig(cnf.MaxSize)
	}
	if cnf.MinSize == 0 {
		return nil
	}
	if cnf.Content.Size < cnf.MinSize {
		return f.Errors.SizeTooSmall(cnf.MinSize)
	}
	return nil
}

func (f Factory) validateWidth(cnf ValidateConfig) *i18np.Error {
	return nil
}

func (f Factory) validateHeight(cnf ValidateConfig) *i18np.Error {
	return nil
}

func (f Factory) GetExtension(content *multipart.FileHeader) string {
	ext := content.Header.Get("Content-Type")
	switch ext {
	case "image/jpeg":
		return "jpg"
	case "image/jpg":
		return "jpg"
	case "image/png":
		return "png"
	case "image/gif":
		return "gif"
	default:
		return ""
	}
}
