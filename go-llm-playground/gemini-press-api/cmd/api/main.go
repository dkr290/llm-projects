package main

import (
	"flag"
	"gemini-press-api/internal/logging"
	"gemini-press-api/internal/server"
	"os"
)

var apiKey string

func main() {
	getEnvs()
	debugFlag := flag.Bool("debugflag", false, "debugflag true or false")
	model := flag.String("model", "gemini-2.0-flash", "default model to use")
	addport := flag.String("addport", "0.0.0.0:9999", "default port and address")
	flag.Parse()

	conf := server.ServerConfig{
		AddrPort:  *addport,
		Model:     *model,
		DebugFlag: *debugFlag,
		ApiKey:    apiKey,
	}
	server.Start(conf)
}

func getEnvs() {
	logger := logging.NewContextLogger("AMLController")

	apiKey = os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		logger.Error().Msg("GEMINI_API_KEY environment variable not set")
		return
	}
}
