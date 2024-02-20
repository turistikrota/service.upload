package event_stream

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/turistikrota/service.upload/src/app/command"
)

func (s Server) ListenUploadPdf(data []byte) {
	cmd := command.UploadPdfCommand{}
	err := json.Unmarshal(data, &cmd)
	if err != nil {
		fmt.Println(err)
		return
	}
	cmd.IsAdmin = true
	_, _ = s.app.Commands.UploadPdf.Handle(context.Background(), cmd)
}

func (s Server) ListenUploadImage(data []byte) {
	cmd := command.UploadImageCommand{}
	err := json.Unmarshal(data, &cmd)
	if err != nil {
		fmt.Println(err)
		return
	}
	cmd.IsAdmin = true
	_, _ = s.app.Commands.UploadImage.Handle(context.Background(), cmd)
}
