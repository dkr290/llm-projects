package server

import (
	"gemini-press-api/internal/handlers"
	"gemini-press-api/internal/logging"
	"net/http"

	"github.com/go-fuego/fuego/option"

	"github.com/go-fuego/fuego"
)

type ServerConfig struct {
	AddrPort  string
	DebugFlag bool
	Model     string
	ApiKey    string
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// allow all options
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Start(c ServerConfig) {
	s := fuego.NewServer(
		fuego.WithAddr(c.AddrPort),
		fuego.WithGlobalMiddlewares(corsMiddleware),
		fuego.WithEngineOptions(fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
			SwaggerURL:       "/docs",
			PrettyFormatJSON: true,
			JSONFilePath:     "docs/openapi.json",
			DisableLocalSave: true,
		}),
		),
	)

	registerRoutes(s, c.DebugFlag, c.ApiKey, c.Model)

	logging.Logger.Info().
		Str("address", c.AddrPort).
		Msg("Server starting")

	if err := s.Run(); err != nil {
		logging.Logger.Fatal().
			Err(err).
			Msg("Server failed")
	}
}

func registerRoutes(s *fuego.Server, debugFlag bool, apiKey, model string) {
	amlController, err := handlers.NewHandler(debugFlag, apiKey, model)
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
