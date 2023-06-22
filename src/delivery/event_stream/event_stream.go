package event_stream

import (
	"api.turistikrota.com/upload/src/app"
	"api.turistikrota.com/upload/src/config"
	"github.com/turistikrota/service.shared/events"
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
	// s.engine.Subscribe(s.Topics.Created, s.ListenEmptyCreated)
}
