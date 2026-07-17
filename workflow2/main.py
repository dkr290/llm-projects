# Run "uv sync" to install the below packages


from openai import OpenAI

client = OpenAI(base_url="http://192.168.1.101:8080/v1/", api_key="not-needed")
EXAMPLES_FILE = Path(__file__).with_name("examples.json")


def load_examples() -> str:
    with EXAMPLES_FILE.open(encoding="utf-8") as file:
        examples = json.load(file)

    return "\n\n".join(
        f"""<example>
<topic>
{example["topic"]}
</topic>
<generated-post>
{example["generated_post"]}
</generated-post>
</example>"""
        for example in examples
    )


def generate_x_post(topic: str) -> str:
    examples = load_examples()
    prompt = f"""
        You are an expert social media manager, and you excel at crafting viral and highly engaging posts for X (formerly Twitter).

        Your task is to generate a post that is concise, impactful, and tailored to the topic provided by the user.
        Avoid using hashtags and lots of emojis (a few emojis are okay, but not too many).

        Keep the post short and focused, structure it in a clean, readable way, using line breaks and empty lines to enhance readability.

        Here's the topic provided by the user for which you need to generate a post:
        <topic>
        {topic}
        </topic>

       Here are some examples of topic and generated posts:
        <examples>
        {examples}
        <examples>
        Please use the tone , language , structure and style of the examples provided above to generate a post that is engaging and relevant to the topic.
"""
    response = client.chat.completions.create(
        model="gemma-4-12b-it-qat-q4_0", messages=[{"role": "user", "content": prompt}]
    )

    return response.choices[0].message.content


def main():
    # user input => AI (LLM) to generate X post => output post

    usr_input = input("What should the post be about? ")
    x_post = generate_x_post(usr_input)
    print("Generated X post")
    print(x_post)


if __name__ == "__main__":
    main()
