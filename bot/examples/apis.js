const axios = require("axios");
const uri="http://129.204.21.76:89/bot"
// const uri="http://127.0.0.1:89/bot"

module.exports.getTask=async function(){
    return axios.get(uri+"/task")
    // axios.get(uri+"/task").then(function(resp){
    //     resp=resp.data
    //     if (resp){
    //         console.debug("get task > ",resp)
    //         if (callback){
    //             callback(resp)
    //         }
    //     }
    // }).catch(function(resp){
    //     console.error(resp)
    // })
}
module.exports.postCmd=function(user,room,content,mentionList){
    let reqt={
        user:user,room:room,content:content,mention:mentionList
    }
    axios.post(uri+"/cmd",reqt).then(function(resp){
        console.debug("post cmd :",resp.data)
    })
}
module.exports.syncContact=function(contacts){
    if (!contacts){
        return
    }
    reqt=new Array()
    for(i=0;i<contacts.length;i++){
        reqt.push({id:contacts[i].id,name:contacts[i].name()})
    }
    console.info("begin sync contacts",reqt.length)
    axios.post(uri+"/sync/user",reqt).then(function(resp){
        console.debug("sync all contacts finished:",reqt.length,resp.data)
    })
}
module.exports.syncAllRoom=function(rooms){
    if (!rooms){
        return
    }
    console.info("begin sync all rooms",rooms.length)
    axios.post(uri+"/sync/room",rooms).then(function(resp){
        console.debug("sync all rooms finished:",resp.data)
    })
}
module.exports.syncUser=function(id,name){
    console.info("begin sync user",id,name)
    axios.post(uri+"/user",{id:id,name:name}).then(function(resp){
        console.debug("sync user:",resp.data)
    })
}
module.exports.syncRoom=function(id,name){
    console.info("begin sync room",id,name)
    axios.post(uri+"/room",{id:id,name:name}).then(function(resp){
        console.debug("sync room",resp.data)
    })
}