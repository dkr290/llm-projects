package llm

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ollama/ollama/api"
)

func LLmreplace(
	valuesFile, chartFile []byte,
	internalRegistry string,
	ollamaHost string,
	model string,
) {
	prompt := fmt.Sprintf(`
Here is a Helm values.yaml file:

%s

And here is the Chart.yaml file:

%s

Please update the values.yaml file based on these rules:
1️⃣ Replace any Docker images from "docker.io" with "%s".
2️⃣ If an image has no tag or tag is null, use "appVersion" from Chart.yaml.
3️⃣ If only "repository: image/name" is given, assume it belongs to "docker.io".
4️⃣ Keep the output **valid YAML**.

Only return the corrected YAML. Do **not** add extra text.
`, string(valuesFile), string(chartFile), internalRegistry)

	fmt.Println(prompt)
	os.Setenv("OLLAMA_HOST", ollamaHost) // Default value
	client, err := api.ClientFromEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 240*time.Second)
	defer cancel()
	genReq := &api.GenerateRequest{
		Model:  model,
		Prompt: prompt,
	}
	// Capture the full response
	var fullResponse string
	genResp := func(r api.GenerateResponse) error {
		if r.Response == "" {
			return fmt.Errorf("empty responce from ollama")
		}
		fullResponse += r.Response
		return nil
	}
	err = client.Generate(ctx, genReq, genResp)
	if err != nil {
		log.Fatal("Error generating response from Ollama:", err)
	}
	// Write modified YAML to a new file
	err = os.WriteFile("modified_values.yaml", []byte(fullResponse), 0644)
	if err != nil {
		log.Fatal("Error writing modified_values.yaml:", err)
	}

	fmt.Println("✅ Updated values.yaml saved as modified_values.yaml!")
}
