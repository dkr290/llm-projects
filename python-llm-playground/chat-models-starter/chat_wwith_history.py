from dotenv import load_dotenv
from langchain_core.messages import HumanMessage, SystemMessage
from langchain_deepseek import ChatDeepSeek

load_dotenv()
llm = ChatDeepSeek(model="deepseek-chat")
msg = [
    SystemMessage("You are an expert in social media content strategy."),
    HumanMessage("Give me a short tip to create engaging content in instagram."),
]

result = llm.invoke(msg)

print(result.content)
