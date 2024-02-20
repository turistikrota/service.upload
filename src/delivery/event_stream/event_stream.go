package event_stream

import (
	"github.com/mixarchitecture/microp/events"
	"github.com/turistikrota/service.upload/src/app"
	"github.com/turistikrota/service.upload/src/config"
)

type Server struct {
	app    app.Application
	Topics config.Topics
	engine events.Engine
}

type Config struct {
	App    app.Application
	Topics config.Topics
	Engine events.Engine
}

func New(config Config) Server {
	return Server{
		app:    config.App,
		engine: config.Engine,
		Topics: config.Topics,
	}
}

func (s Server) Load() {
	s.engine.Subscribe(s.Topics.Upload.UploadPDF, s.ListenUploadPdf)
	s.engine.Subscribe(s.Topics.Upload.UploadImage, s.ListenUploadImage)
}
