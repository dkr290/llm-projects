import gradio as gr
import requests

def query_golang_app(text):
    # Send request to Golang app
    url = "http://localhost:3000/generate"
    payload = {"prompt": text}
    
    try:
        response = requests.post(url, json=payload)

        if response.ok:  # Checks for any 2xx status code
            return response.json().get("generated_text", "No correction generated.")
        else:
            return f"Error: {response.status_code} - {response.text}"

    except requests.exceptions.RequestException as e:
        return f"Request failed: {e}"
# Define Gradio interface
interface = gr.Interface(
    fn=query_golang_app,
    inputs=gr.Textbox(lines=5,placeholder="Enter text with grammer or spelling mistakes"),
    outputs=gr.Textbox(label="Corrected text"),
    title="AI powered Grammer and Spell Checker",
    description="Enter a prompt to get a response from Ollama via a Golang backend."
)

if __name__ == "__main__":
    interface.launch(server_name="0.0.0.0")
