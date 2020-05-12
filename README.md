# repe-wechat-assistant 
[![Powered by Wechaty](https://img.shields.io/badge/Powered%20By-Wechaty-green.svg)](https://github.com/chatie/wechaty)
[![Wechaty开源激励计划](https://img.shields.io/badge/Wechaty-开源激励计划-green.svg)](https://github.com/juzibot/Welcome/wiki/Everything-about-Wechaty)

repe is a wechat group assistant,help administrator manage WeChat group.
## fefeature 
- sync user and room to service automatically 
- response WeChat message immediately 
- support dynamic command configuration 
## structure 
### /backend 
Writed by golang,is the service apply api for bot
### /bot
Writed by nodejs,is bot self,only do the bot things,all command will be request to service and get response from service.
### get-start 
#### start service 
<pre>
cd backend/main
go build 
./main
</pre>
#### start bot
<pre>
cd bot
npm start
</pre>
### test UI 
- After start backend service,you can open http://localhost/example/actv with browser,add new activity to a WeChat group or user. The command filed you set,is the way how user join that activity.
- Also there is a bot simulater at http://localhost:89/example/1 , use it to test when bot can't be setuped. 
