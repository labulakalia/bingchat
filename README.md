
# Bing Ai Chat Copilot By Golang 

## Requirements
A Microsoft Account with early access to https://bing.com/chat (Required)
Required in a supported country with New Bing Or set http_proxy

## Install
- Install [Cookie-Editor](https://chrome.google.com/webstore/detail/cookie-editor/hlkenndednhfkekhgcdicdfddnkalmdm?hl=en),
- Export `bing.com` cookies to save json file
- Download `bingchat` for your platform

## Usage
```shell
➜  cmd ./bingchat -c cookies.json
local http proxy is set
Bing Ai Chat Copilot
/reset [styles]
      1  more create
      2  more balance
      3  more precise
/quit
  quit
/help
  print help
Current Style: Balance 
Ask> hello
Hello! How can I help you today? 😊
Prompt suggest
1: What is the weather like today?
2: What is the latest news?
3: Can you tell me a joke?
Ask> 2
Here are some of the latest news sources you can check out:
- [The Guardian](https://www.theguardian.com/world)
- [NBC News](https://www.nbcnews.com/latest-stories)
- [RTÉ News](https://www.rte.ie/news/)

Which one would you like to check out?
Prompt suggest
1: What is the latest news on politics?
2: What is the latest news on sports?
3: What is the latest news on entertainment?
Ask> /reset 1
Switch Style Create
Ask> 
```