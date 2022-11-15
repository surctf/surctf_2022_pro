import requests as reqs

url = "https://api.telegram.org/bot5631770980:AAExlmzAZ1fbMuWwHd5Oa7oHmrtnvdRMuB8/"

resp = reqs.get(url + "getMe")
print(resp.text)

resp = reqs.get(url + "getChatAdministrators", params={"chat_id":-875674057})
print(resp.text)