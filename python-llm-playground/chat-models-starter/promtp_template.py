from dotenv import load_dotenv
from langchain_core.prompts import ChatPromptTemplate
from langchain_deepseek import ChatDeepSeek

load_dotenv()
llm = ChatDeepSeek(model="deepseek-chat")


# template = "Write a {tone} email to {company} expressing interest in the {position} position, mentioning {skill} as a key strength. Keep it to 4 lines max"

# prompt_template = ChatPromptTemplate.from_template(template)

# prompt =  prompt_template.invoke({
#     "tone": "energetic",
#     "company": "samsung",
#     "position": "AI Engineer",
#     "skill": "AI"
# })

messages = [
    ("system", "You are a comedian who tells jokes about {topic}."),
    ("human", "Tell me {joke_count} jokes."),
]

prompt_template = ChatPromptTemplate.from_messages(messages)
prompt = prompt_template.invoke({"topic": "lawyers", "joke_count": 3})
result = llm.invoke(prompt)
print(result)
