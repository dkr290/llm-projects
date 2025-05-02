package main

import (
	"flag"
	"gemini-press-api/internal/logging"
	"gemini-press-api/internal/server"
	"os"
	"strings"
)

var (
	apiKey    string
	appURL    string
	modelList []string
)

func main() {
	debugFlag := flag.Bool("debugflag", false, "debugflag true or false")
	model := flag.String("model", "gemini-2.0-flash", "default model to use")
	addport := flag.String("addport", "9999", "default port and address")
	public_url := flag.String(
		"publicurl",
		"https://example.com",
		"The Public url to use only for production deployments",
	)
	flag.Parse()
	getEnvs()

	if os.Getenv("PUBLIC_URL") != "" {
		appURL = os.Getenv("PUBLIC_URL")
	} else {
		appURL = *public_url
	}
	modelList = getModelList()
	if modelList == nil {
		modelList = append(modelList, *model)
	}
	conf := server.ServerConfig{
		AddrPort:  *addport,
		Models:    modelList,
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
}

func getModelList() []string {
	modelList := os.Getenv("MODEL_LIST")
	var trimmedModels []string
	if modelList != "" {
		models := strings.Split(modelList, ",")
		for _, model := range models {
			trimmedModel := strings.TrimSpace(model)
			if trimmedModel != "" {
				trimmedModels = append(trimmedModels, trimmedModel)
			}
		}
	} else {
		trimmedModels = nil
	}
	return trimmedModels
}
