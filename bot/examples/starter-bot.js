const {
  Wechaty,
  ScanStatus,
  log,
  Friendship,
}               = require('wechaty')
const { PuppetPadplus } =require('wechaty-puppet-padplus')
const apis=require("./apis.js")
const WECHATY_PUPPET_PADPRO_TOKEN = 'puppet_padplus_898ac68f127f21aa'

const puppet = new PuppetPadplus({
  token: WECHATY_PUPPET_PADPRO_TOKEN,
})

const bot = new Wechaty({
  puppet,
})
bot.on('scan',    onScan) //扫码触发
bot.on('login',   onLogin) //登录成功
bot.on('logout',  onLogout) //登出
bot.on('message', onMessage) //收到消息
bot.on('friendship',onFriendship) //新联系人 
bot.on('room-join',onOneJoinRoom) //新人加入群聊
bot.on('room-invite',onRoomInvite) //邀请加入群聊

function onScan (qrcode, status) {
  if (status === ScanStatus.Waiting || status === ScanStatus.Timeout) {
    require('qrcode-terminal').generate(qrcode, { small: true })  // show qrcode on console

    const qrcodeImageUrl = [
      'https://api.qrserver.com/v1/create-qr-code/?data=',
      encodeURIComponent(qrcode),
    ].join('')

    log.info('StarterBot', 'onScan: %s(%s) - %s', ScanStatus[status], status, qrcodeImageUrl)

  } else {
    log.info('StarterBot', 'onScan: %s(%s)', ScanStatus[status], status)
  }
}
// 联系人缓存
var allContact ={}
// 群缓存
var allRooms={}

function syncUsers(contacts){
  if (Array.isArray(contacts)){
    for(i=0;i<contacts.length;i++){
      allContact[contacts[i].id]=contacts[i]
    }
    apis.syncContact(contacts)
  }else{
    allContact[contacts.id]=contacts
    apis.syncUser(contacts.id,contacts.name())
  }
}
async function syncRooms(rooms){
  if (Array.isArray(rooms)){
    for(i=0;i<rooms.length;i++){
      allRooms[rooms[i].id]=rooms[i]
    }
    rms=await convertRoom(rooms)
    apis.syncAllRoom(rms)
  }else{
    allRooms[rooms.id]=rooms
    rm=await convertRoom(rooms)
    apis.syncRoom(rm.id,rm.name)
  }
}
async function convertRoom(room){
  if (Array.isArray(room)){
    result=new Array()
    for(i=0;i<room.length;i++){
      result.push({id:room[i].id,name:await room[i].topic()})
    }
    return result
  }else{
    return{
      id:room.id,
      name:await room.topic()
    }
  }
}
async function onLogin (user) {
  console.log(`${user} login`)
  // 同步联系人信息
  const contacts = await bot.Contact.findAll()
  syncUsers(contacts)

  // 同步群聊信息
  let rooms = await bot.Room.findAll()
  
  await syncRooms(rooms)
  for(i=0;i<rooms.length;i++){
    rooms[i].on('topic',onRoomTopic)
    let members=await rooms[i].memberAll()
    // 同步群内成员信息
    syncUsers(members)
  }

  // 定时拉取待发送消息
  setInterval(function(){
    apis.getTask().then(autoSendMsg)
  },1000)
  console.log("----------setup bot successful-----------")
}
async function autoSendMsg(resp){
    resp=resp.data
    if(!resp||(!resp.user&&!resp.room)){
      return
    }
    console.debug("receive new task",resp)
    // 私聊
    if(resp.user&&!resp.room){
      user=allContact[resp.user]
      if (!user){
        return
      }
      await user.say(resp.content)
      return
    }
    //群聊
    if(resp.room){
      room=allRooms[resp.room]
      if (!room){
        return
      }
      //是否@user
      if (resp.user){
        await room.say(resp.content,allContact[resp.user])
      }else{
        await room.say(resp.content)
      }
    }
}
async function onRoomTopic(room, topic, oldTopic, changer){
  console.log(`Room topic changed from ${oldTopic} to ${topic} by ${changer.name()}`)
  syncRooms(room)
}
function onLogout(user) {
  console.log(`${user} logout`)
}

async function onMessage (msg) {
  console.log(msg)
  const contact = msg.from()
  if (contact.self()){
    return
  }
  const text = msg.text()
  const room = msg.room()
  if (!text||!contact){
    return
  }
  let roomID =""
  if (room){
    roomID=room.id
  }
  let contactList = await msg.mentionList()
  let ids=[]
  for(i=0;i<contactList.length;i++){
    ids.push(contactList[i].id)
  }
  apis.postCmd(contact.id,roomID,text,ids)
}
async function onFriendship(friendship){
  if(friendship.type() === Friendship.Type.Receive){
    // 1. receive new friendship request from new contact    
    const contact = friendship.contact()    
    let result = await friendship.accept()      
    if(result){
      console.log(`Request from ${contact.name()} is accept succesfully!`)
    } else {        
      console.error(`Request from ${contact.name()} failed to accept!`)
      return   
    }      
  } else if (friendship.type() === Friendship.Type.Confirm) { 
    // 2. confirm friendship      
      console.log(`new friendship confirmed with ${friendship.contact().name()}`)   
    }
    // 同步这个联系人信息
    user=friendship.contact()
    syncUsers(user)
}
async function onOneJoinRoom(room, inviteeList, inviter){
  if (!inviteeList){
    return
  }
  console.debug(inviteeList)
  console.log(`Room ${await room.topic()} got ${inviteeList.length} members, invited by ${inviter}`)
  for(i=0;i<inviteeList.length;i++){
    if (inviteeList[i].self()){
      // 如果是自己被邀请，同步群信息
      room.on('topic',onRoomTopic)
      await syncRooms(room)
      syncUsers(await room.memberAll())
    }else{
      // 同步用户信息
      syncUsers(inviteeList[i])
    }
  }
}
async function onRoomInvite(roomInvitation){
  try {    
    console.log(`received room-invite event.`)
    // 自动接收群邀请
    await roomInvitation.accept()

  } catch (e) {    
    console.error(e)  
  }

}
async function onStart(){
  console.log("bot start successly")
}
bot.start()
.then(onStart)
.catch(e => console.error(e))
