from dotenv import load_dotenv
from langchain_core.messages import AIMessage, HumanMessage, SystemMessage
from langchain_deepseek import ChatDeepSeek

load_dotenv()
llm = ChatDeepSeek(model="deepseek-chat")

system_message = SystemMessage(content="You are helpful AI assistant")
chat_history = []
chat_history.append(system_message)

while True:
    query = input("You:")
    if query.lower() == "exit":
        break
    chat_history.append(HumanMessage(content=query))  # adding user message
    result = llm.invoke(chat_history)
    resp = result.content
    chat_history.append(AIMessage(content=resp))

    print(f"AI: {resp}")


print("----Message History----")
print(chat_history)
