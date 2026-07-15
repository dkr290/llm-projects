# Run "uv sync" to install the below packages

import requests


def generate_x_post(topic: str) -> str:
    prompt = f"""
        You are an expert social media manager, and you excel at crafting viral and highly engaging posts for X (formerly Twitter).

        Your task is to generate a post that is concise, impactful, and tailored to the topic provided by the user.
        Avoid using hashtags and lots of emojis (a few emojis are okay, but not too many).

        Keep the post short and focused, structure it in a clean, readable way, using line breaks and empty lines to enhance readability.

        Here's the topic provided by the user for which you need to generate a post:
        <topic>
        {topic}
        </topic>
"""
    payload = {
        "model": "gemma-4-12b-it-qat-q4_0",
        "messages": [{"role": "user", "content": prompt}],
    }
    response = requests.post(
        "http://192.168.1.101:8080/v1/chat/completions",
        json=payload,
        headers={
            "Content-Type": "application/json",
        },
        timeout=60,
    )
    response.raise_for_status()

    response_text = (
        response.json().get("choices", [{}])[0].get("message", {}).get("content", "")
    )

    return response_text


def main():
    # user input => AI (LLM) to generate X post => output post

    usr_input = input("What should the post be about? ")
    x_post = generate_x_post(usr_input)
    print("Generated X post")
    print(x_post)


if __name__ == "__main__":
    main()
