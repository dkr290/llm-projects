import gradio as gr
import requests

# Legal Document Templates
LEGAL_TEMPLATES = {
    "rental agreement": "Generate a rental agreement between {party1} (tenant) and {party2} (landlord) for {duration} months.",
    "employment contract": "Generate an employment contract between {party1} (employee) and {party2} (employer) with a salary of {salary} per year.",
    "business partnership agreement": "Draft a business partnership agreement between {party1} and {party2}, defining responsibilities and profit-sharing terms.",
    "nda": "Generate a non-disclosure agreement (NDA) between {party1} and {party2} to protect confidential business information.",
}


def query_golang_app(doc_type, party1, party2, duration="", salary=""):
    # Send request to Golang app
    url = "http://localhost:3000/generate"

    if doc_type not in LEGAL_TEMPLATES:
        return "Invalid document type. Please choose from rental agreement, employment contract, business partnership agreement, or NDA."

    text = LEGAL_TEMPLATES[doc_type].format(
        party1=party1, party2=party2, duration=duration, salary=salary
    )

    payload = {"prompt": text}

    try:
        response = requests.post(url, json=payload, timeout=50)
        if response.ok:  # Checks for any 2xx status code
            return response.json().get("generated_text", "No document generated.")

        return f"Error: {response.status_code} - {response.text}"

    except requests.exceptions.RequestException as e:
        return f"Request failed: {e}"


# Define Gradio interface
interface = gr.Interface(
    fn=query_golang_app,
    inputs=[
        gr.Radio(
            [
                "rental agreement",
                "employment contract",
                "business partnership agreement",
                "nda",
            ],
            label="Document Type",
        ),
        gr.Textbox(label="party 1 Name"),
        gr.Textbox(label="party 2 Name"),
        gr.Textbox(label="Duration (if applicable, in months)", placeholder="e.g., 12"),
        gr.Textbox(
            label="Salary (if applicable, per year)", placeholder="e.g., $50,000"
        ),
    ],
    outputs=gr.Textbox(label="Generated Legal Document"),
    title="AI powered Legal Assistant",
    description="Select a document type, enter party names, and generate a professional legal contract.",
)

if __name__ == "__main__":
    interface.launch(server_name="0.0.0.0")
