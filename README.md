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
但是如果采用这种方案，那么客户端和服务端都需要至少有两个端口来负责简单的痛惜，显然是不现实的，所以这种方案不可取
### 2.使用收发服务器
使用两个服务端，一个进行收取消息，另一个负责把消息广播 ，这个方案也要要求服务器端两个，十分麻烦
### 3 使用Kratos微服务框架 
<br>详情：https://juejin.cn/post/7104264158467588133#heading-1
## 最终采用的方案
<br>最终采用了使用mysql+grpc的通信方案，服务器端和客户端通过同时访问mysql来实现，服务器统一将客户端的信息输入数据库中，实现通信<br><br>
<br>服务定义为<br>

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
  string time = 4;
}
message MessageNum{
  int64 messnum = 1;
}
service Chat{
  rpc SendAll(Message )returns(Message);
  rpc GetMessNum(UserId)returns(MessageNum);
}

 ```
 
 <br>运行结果:
 <br>
 <br>
<img width="574" alt="image" src="https://user-images.githubusercontent.com/96430610/198608106-8a562a04-e1ae-4de8-ae70-fa11d1a3e60d.png">
<br>
<br>
虽然在这种单机版中，grpc的作用被明显弱化了，但是当此服务部署到web上时，grpc的作用就凸现出来了，所有客户把消息通过grpc统一交给服务端管理，不仅能节省很多时间，还方便管理，而服务端的总消息库又是对所有用户可见的，这也就实现了广播。
## 在web中使用
gRPC由于自身的特性是不适合使用对外部服务的，所以如果要想在一个web项目中使用这个基于gRPC的聊天系统的话，那么使用的方案如下：
### gRPC是一种进程间的通信方式，所以以上的client模块和server模块都是放在服务器端运行的，而每个client模块则在与server端对接传输数据的同时，负责和web请求对接，每个用户在聊天的时候对应着一个client，这样的话就可以在服务器进程中实现通信。<br>所以只需要把client和server对外部开放一下就可以
