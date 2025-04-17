package main

import (
	"flag"
	"gemini-press-api/internal/logging"
	"gemini-press-api/internal/server"
	"os"
)

var (
	apiKey string
	appURL string
)

func main() {
	getEnvs()
	debugFlag := flag.Bool("debugflag", false, "debugflag true or false")
	model := flag.String("model", "gemini-2.0-flash", "default model to use")
	addport := flag.String("addport", "9999", "default port and address")
	flag.Parse()

	conf := server.ServerConfig{
		AddrPort:  *addport,
		Model:     *model,
		DebugFlag: *debugFlag,
		ApiKey:    apiKey,
		PublicURL: appURL,
	}
	server.Start(conf)
}

func getEnvs() {
	logger := logging.NewContextLogger("AMLController")

	apiKey = os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		logger.Error().Msg("GEMINI_API_KEY environment variable not set")
		os.Exit(1)
	}
	appURL = os.Getenv("PUBLIC_URL")
	if appURL == "" {
		logger.Error().Msg("We need production URL PUBLIC_URL to be set")
		os.Exit(1)
	}
}
