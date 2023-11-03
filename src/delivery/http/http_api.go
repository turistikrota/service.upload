package http

import (
	"mime/multipart"

	"github.com/gofiber/fiber/v2"
	httpI18n "github.com/mixarchitecture/microp/server/http/i18n"
	"github.com/mixarchitecture/microp/server/http/result"
	"github.com/turistikrota/service.shared/server/http/auth/current_account"
	"github.com/turistikrota/service.upload/src/app/command"
	"github.com/turistikrota/service.upload/src/delivery/http/dto"
	"github.com/turistikrota/service.upload/src/domain/cdn"
)

type fileRequest struct {
	FileName    string
	DirName     string
	RandomName  bool
	Slugify     bool
	IsAdmin     bool
	MinifyLevel cdn.MinifyLevel
	Content     *multipart.FileHeader
}

func (h Server) UploadImage(ctx *fiber.Ctx) error {
	file, err := h.validateAdmin(ctx, Fields.Image, Messages.Error.ImageNotFound)
	if err != nil {
		return err
	}
	res, error := h.app.Commands.UploadImage.Handle(ctx.UserContext(), command.UploadImageCommand{
		RandomName:  file.RandomName,
		Content:     file.Content,
		IsAdmin:     file.IsAdmin,
		FileName:    file.FileName,
		Dir:         file.DirName,
		MinifyLevel: file.MinifyLevel,
		Slugify:     file.Slugify,
	})
	return result.IfSuccessDetail(error, ctx, h.i18n, Messages.Success.ImageUploaded, func() interface{} {
		return dto.Response.ImageUploaded(res)
	})
}

func (h Server) UploadSvg(ctx *fiber.Ctx) error {
	file, err := h.validateAdmin(ctx, Fields.Svg, Messages.Error.ImageNotFound)
	if err != nil {
		return err
	}
	res, error := h.app.Commands.UploadSvg.Handle(ctx.UserContext(), command.UploadSvgCommand{
		RandomName: file.RandomName,
		Content:    file.Content,
		IsAdmin:    file.IsAdmin,
		FileName:   file.FileName,
		Dir:        file.DirName,
	})
	return result.IfSuccessDetail(error, ctx, h.i18n, Messages.Success.ImageUploaded, func() interface{} {
		return dto.Response.SvgUploaded(res)
	})
}

func (h Server) UploadPdf(ctx *fiber.Ctx) error {
	file, err := h.validateAdmin(ctx, Fields.Pdf, Messages.Error.PdfNotFound)
	if err != nil {
		return err
	}
	res, error := h.app.Commands.UploadPdf.Handle(ctx.UserContext(), command.UploadPdfCommand{
		RandomName: file.RandomName,
		Content:    file.Content,
		IsAdmin:    file.IsAdmin,
		FileName:   file.FileName,
		Dir:        file.DirName,
	})
	return result.IfSuccessDetail(error, ctx, h.i18n, Messages.Success.PdfUploaded, func() interface{} {
		return dto.Response.PdfUploaded(res)
	})
}

func (h Server) UploadMarkdown(ctx *fiber.Ctx) error {
	file, err := h.validateAdmin(ctx, Fields.Markdown, Messages.Error.MarkdownNotFound)
	if err != nil {
		return err
	}
	res, error := h.app.Commands.UploadMarkdown.Handle(ctx.UserContext(), command.UploadMarkdownCommand{
		RandomName: file.RandomName,
		Content:    file.Content,
		IsAdmin:    file.IsAdmin,
		FileName:   file.FileName,
		Dir:        file.DirName,
	})
	return result.IfSuccessDetail(error, ctx, h.i18n, Messages.Success.MarkdownUploaded, func() interface{} {
		return dto.Response.MarkdownUploaded(res)
	})
}

func (h Server) UploadAvatar(ctx *fiber.Ctx) error {
	avatar, err := h.validateAvatar(ctx)
	if err != nil {
		return err
	}
	name := current_account.Parse(ctx)
	res, error := h.app.Commands.UploadAvatar.Handle(ctx.UserContext(), command.UploadAvatarCommand{
		Content:  avatar,
		UserName: name,
	})
	return result.IfSuccessDetail(error, ctx, h.i18n, Messages.Success.AvatarUploaded, func() interface{} {
		return dto.Response.AvatarUploaded(res)
	})
}

func (h Server) validateAdmin(ctx *fiber.Ctx, field string, errorMsg string) (*fileRequest, error) {
	image, err := ctx.FormFile(field)
	if err != nil {
		l, a := httpI18n.GetLanguagesInContext(h.i18n, ctx)
		return nil, result.Error(h.i18n.Translate(errorMsg, l, a))
	}
	fileName := ctx.FormValue("fileName", "")
	dirName := ctx.FormValue("dirName", "")
	randomName := ctx.FormValue("randomName", "true")
	slugify := ctx.FormValue("slugify", "false")
	minifyLevel := ctx.FormValue("minifyLevel", "0")
	level := cdn.MinifyLevelLow
	if minifyLevel == "medium" {
		level = cdn.MinifyLevelMedium
	} else if minifyLevel == "high" {
		level = cdn.MinifyLevelHigh
	} else if minifyLevel == "none" {
		level = cdn.MinifyLevelNone
	}
	u := ctx.Locals(field)
	return &fileRequest{
		FileName:    fileName,
		Slugify:     slugify == "true",
		DirName:     dirName,
		RandomName:  randomName == "true",
		IsAdmin:     u != nil && u.(string) == "true",
		Content:     image,
		MinifyLevel: level,
	}, nil
}

func (h Server) validateAvatar(ctx *fiber.Ctx) (*multipart.FileHeader, error) {
	image, err := ctx.FormFile("avatar")
	if err != nil {
		l, a := httpI18n.GetLanguagesInContext(h.i18n, ctx)
		return nil, result.Error(h.i18n.Translate(Messages.Error.AvatarNotFound, l, a))
	}
	return image, nil
}
