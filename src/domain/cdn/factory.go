package cdn

import (
	"bytes"
	"image"
	"image/png"
	"mime/multipart"

	"github.com/disintegration/imaging"
	"github.com/google/uuid"
	"github.com/mixarchitecture/i18np"
	"github.com/nfnt/resize"
)

type Factory struct {
	Errors    Errors
	watermark image.Image
}

func NewFactory() Factory {
	return Factory{
		Errors:    newCdnErrors(),
		watermark: loadWatermark(),
	}
}

func loadWatermark() image.Image {
	watermarkImage, err := imaging.Open("assets/watermark.png")
	if err != nil {
		panic(err)
	}
	return watermarkImage
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

func (f Factory) NewImage(cnf ValidateConfig) ([]byte, *i18np.Error) {
	err := f.Validate(cnf)
	if err != nil {
		return nil, err
	}
	file, error := cnf.Content.Open()
	if error != nil {
		return nil, f.Errors.InternalError()
	}
	defer file.Close()
	bytes, _err := f.watermarkFromMultipart(file)
	if _err != nil {
		return nil, f.Errors.InternalError()
	}
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

func (f Factory) watermarkImage(originalImage image.Image) (image.Image, *i18np.Error) {
	originalWidth := originalImage.Bounds().Dx()
	originalHeight := originalImage.Bounds().Dy()
	scaledWatermark := imaging.Resize(f.watermark, originalWidth, 0, imaging.Lanczos)
	result := imaging.Overlay(originalImage, scaledWatermark, image.Pt((originalWidth-scaledWatermark.Bounds().Dx())/2, (originalHeight-scaledWatermark.Bounds().Dy())/2), 1.0)
	return result, nil
}

func (f Factory) minifyImage(originalImage image.Image) image.Image {
	const maxWidth = 1920
	const maxHeight = 1080
	width := originalImage.Bounds().Dx()
	height := originalImage.Bounds().Dy()
	maxSizeInBytes := int64(0.5 * 1024 * 1024) // 0.5 MB

	if int64(width*height*3) <= maxSizeInBytes {
		return originalImage
	}

	resizedImage := resize.Resize(uint(maxWidth), uint(maxHeight), originalImage, resize.Lanczos3)

	for int64(resizedImage.Bounds().Dx()*resizedImage.Bounds().Dy()*3) > maxSizeInBytes {
		resizedImage = resize.Resize(uint(width/2), uint(height/2), resizedImage, resize.Lanczos3)
		width = width / 2
		height = height / 2
	}

	return resizedImage
}

func (f Factory) watermarkFromMultipart(file multipart.File) ([]byte, *i18np.Error) {
	originalImage, _, err := image.Decode(file)
	if err != nil {
		return nil, f.Errors.InternalError()
	}
	watermarkedImage, error := f.watermarkImage(originalImage)
	if error != nil {
		return nil, f.Errors.InternalError()
	}
	watermarkedImage = f.minifyImage(watermarkedImage)
	return f.watermarkToBytes(watermarkedImage)
}

func (f Factory) watermarkToBytes(img image.Image) ([]byte, *i18np.Error) {
	buffer := new(bytes.Buffer)
	err := png.Encode(buffer, img)
	if err != nil {
		return nil, f.Errors.InternalError()
	}
	return buffer.Bytes(), nil
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
	case "text/markdown":
		return "md"
	case "text/plain":
		return "txt"
	default:
		return ""
	}
}
