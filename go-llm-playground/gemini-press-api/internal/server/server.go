package server

import (
	"fmt"
	"gemini-press-api/internal/handlers"
	"gemini-press-api/internal/logging"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-fuego/fuego/option"
	"github.com/rs/cors"

	"github.com/go-fuego/fuego"
)

type ServerConfig struct {
	AddrPort  string
	DebugFlag bool
	Models    []string
	ApiKey    string
	PublicURL string
}

// func corsMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		// allow all options
// 		w.Header().Set("Access-Control-Allow-Origin", "*")
// 		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
// 		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
//
// 		if r.Method == http.MethodOptions {
// 			w.WriteHeader(http.StatusOK)
// 			return
// 		}
//
// 		next.ServeHTTP(w, r)
// 	})
// }

// can be done also with this package the above will be kept commented out just for the example
func Start(c ServerConfig) {
	cs := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
		Debug:            false,
	})
	s := fuego.NewServer(
		fuego.WithAddr(fmt.Sprintf(":%s", c.AddrPort)),
		fuego.WithGlobalMiddlewares(cs.Handler),
		fuego.WithEngineOptions(fuego.WithOpenAPIConfig(fuego.OpenAPIConfig{
			SwaggerURL:       "/docs",
			PrettyFormatJSON: true,
			JSONFilePath:     "docs/openapi.json",
			DisableLocalSave: false,
		}),
		),
	)
	s.OpenAPI.Description().Servers = append(s.OpenAPI.Description().Servers, &openapi3.Server{
		URL:         c.PublicURL,
		Description: "Production",
	})
	s.OpenAPI.Description().Servers = append(s.OpenAPI.Description().Servers, &openapi3.Server{
		URL:         fmt.Sprintf("%s:%s", "http://localhost", c.AddrPort),
		Description: "localhost",
	})

	registerRoutes(s, c.DebugFlag, c.ApiKey, c.Models)

	logging.Logger.Info().
		Str("address", c.AddrPort).
		Msg("Server starting")

	if err := s.Run(); err != nil {
		logging.Logger.Fatal().
			Err(err).
			Msg("Server failed")
	}
}

func registerRoutes(s *fuego.Server, debugFlag bool, apiKey string, models []string) {
	amlController, err := handlers.NewHandler(debugFlag, apiKey, models)
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
