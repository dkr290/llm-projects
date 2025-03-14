import gradio as gr
import requests

def query_golang_app(prompt, model="llama2"):
    # Send request to Golang app
    url = "http://localhost:8080/api/ollama"
    payload = {"prompt": prompt, "model": model}
    response = requests.post(url, json=payload)
    
    if response.status_code == 200:
        return response.json()["response"]
    else:
        return "Error: Could not get response from Golang app"

# Define Gradio interface
interface = gr.Interface(
    fn=query_golang_app,
    inputs=["text", gr.Dropdown(choices=["llama2", "mistral"], label="Model")],
    outputs="text",
    title="Ollama Chat via Golang",
    description="Enter a prompt and select a model to get a response from Ollama via a Golang backend."
)

# Launch the interface
interface.launch()
