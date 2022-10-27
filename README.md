# 使用原生go环境搭建基于gRPC的广播系统
## 与socket的不同
<br>  这里的广播不像之前的socket编程一样简单，主要的差别在于gRPC是一个通信协议，而socket更像是一个文件，只要建立连接后，谁都可以在里面读和写，所以server端和client端之间的交流相对来说比较简单，而gRPC的不同之处在于，gRPC更像是一个通信协议，每次使用gRPC通信都需要得到通信双方的认可<br>
正是由于这个特性，导致使用gRPC进行广播的时候就很难确定出什么时候client端该向server端索要消息数据了。
<br>同时，由服务定义生成的server骨架和client存根也导致了gRPC开发中proto文件十分重要，具有牵一发而动全身的影响，所以，需要注意的是，在开发gRPC应用之前，一定要确定好proto文件！
对于以上的通信问题，有几种可能可行的办法：
### 1.使用NotifyManager
<br>NotifyManager，本质上也是一个服务，client端在会与用于发送消息的服务和NotifyManager同时消息流连接，同时，NM和消息服务器相连，一旦有客户端向消息服务器发送消息，NM给客户发的消息流就会发生改变，这也就是提醒用户此时可以向服务端发送请求了，此时的服务定义应为：

```

syntax="proto3";
option go_package="github.com/Sean04111/chat_gRPC";
package chat;
message User{
  int64 id = 1;
  string name = 2;
}
message UserId{
  int64 id =1;
}
message Message{
  int64 id = 1;
  string speakername=2;
  string content =3;
}
message NotiCode{
  int64 code = 1;
}
service Chat{
  rpc SendAll(Message )returns(Message);
  rpc HaveAll(UserId)returns(Message);
}
service Notify{
  rpc Notify(UserId)returns(stream NotiCode);
  
  
  ```


<br>之间的通信可以通过并发管道实现
### 2.使用收发服务器
使用两个服务端，一个进行收取消息，另一个负责把消息广播 （目前还没有实践过）
### 3 使用Kratos微服务框架 
<br>详情：https://juejin.cn/post/7104264158467588133#heading-1
## 已经完成的工作
<br>目前基于客户端流RPC模式实现了client端向服务器端发送消息流，并且server端能够正常接受，可以选择存全局变量里，也可以选择讯数据库里<br><br>
![K}M`5H(WO}_@LSMNDM%8YCP](https://user-images.githubusercontent.com/96430610/198402595-0c184612-1061-4d13-9a6f-8dc1892e92f3.png)
同时方案一还在实验中......

