package server

import (
	"flag"
	"press-detective/internal/controllers"
	"press-detective/internal/logging"

	"github.com/go-fuego/fuego/option"

	"github.com/go-fuego/fuego"
)

func Start() {
	s := createServer()
	registerRoutes(s)

	logging.Logger.Info().
		Str("address", "localhost:9999").
		Msg("Server starting")

	if err := s.Run(); err != nil {
		logging.Logger.Fatal().
			Err(err).
			Msg("Server failed")
	}
}

func createServer() *fuego.Server {
	return fuego.NewServer(
		fuego.WithEngineOptions(fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
			SwaggerURL:   "/docs",
			SpecURL:      "/docs/openapi.json", // URL to serve the openapi json spec
			JSONFilePath: "doc/openapi.json",   // Local path to save the openapi json spec

		})),
	)
}

func registerRoutes(s *fuego.Server) {
	debugFlag := flag.Bool("debugflag", false, "debugflag true or false")
	flag.Parse()
	amlController, err := controllers.NewAMLController(*debugFlag)
	if err != nil {
		logging.Logger.Fatal().
			Err(err).
			Msg("Failed to create AML controller")
	}
	fuego.Post(
		s,
		"/search",
		amlController.SearchHandler,
		option.Summary("Search for target name"),
	)
	fuego.Get(s, "/", amlController.RootHandler, option.Summary("Home page"))
	fuego.Get(s, "/ping", amlController.PingHandler, option.Summary("Health check"))
}
