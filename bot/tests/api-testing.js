const apis=require("../examples/apis.js")

async function main(){
    setInterval(function(){
        apis.getTask(function(resp){
            console.log("robot>>>>",resp)
        })
    },1000)
    apis.syncUser("wx_uu1","臭弟弟")
    apis.syncRoom("wx_rr1","主动出击小组")
    apis.postCmd("wx_uu1","wx_rr1","签到")
    apis.postCmd("wx_uu1","wx_rr1","帮助")
    apis.postCmd("wx_uu1","wx_rr1","余额")
}

main()
.then()
.catch(e => {
  console.error(e)
  process.exit(1)
})