package main

import (
	"flag"
	"gemini-press-api/internal/logging"
	"gemini-press-api/internal/server"
	"os"
	"strings"
)

var (
	apiKeys   []string
	appURL    string
	modelList []string
)

func main() {
	debugFlag := flag.Bool("debugflag", false, "debugflag true or false")
	model := flag.String("model", "gemini-2.0-flash", "default model to use")
	addport := flag.String("addport", "9999", "default port and address")
	publicURL := flag.String(
		"publicurl",
		"https://example.com",
		"The Public url to use only for production deployments",
	)
	flag.Parse()
	logging.Init(*debugFlag)

	apiKeys = getEnvs()

	if os.Getenv("PUBLIC_URL") != "" {
		appURL = os.Getenv("PUBLIC_URL")
	} else {
		appURL = *publicURL
	}
	modelList = getModelList()
	if modelList == nil {
		modelList = append(modelList, *model)
	}
	conf := server.ServerConfig{
		AddrPort:  *addport,
		Models:    modelList,
		DebugFlag: *debugFlag,
		ApiKeys:   apiKeys,
		PublicURL: appURL,
	}
	server.Start(conf)
}

func getEnvs() []string {
	logger := logging.NewContextLogger("AMLController")

	apiKeys := os.Getenv("GEMINI_API_KEY")
	apiKeys = strings.TrimSpace(apiKeys)
	keyList := strings.Split(apiKeys, ",")
	if keyList[0] == "" {
		logger.Error().Msg("GEMINI_API_KEY environment variable not set")
		os.Exit(1)
	}
	return keyList
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
