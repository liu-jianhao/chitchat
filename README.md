# ChitChat

ChitChat的数据模型：
1. User——表示论坛的用户信息
2. Session——表示论坛用户当前的登录会话
3. Thread——表示论坛里面的帖子，每一个帖子都记录了多个论坛用户之间的对话
4. Post——表示用户在帖子里面添加的回复

## 注意
如果直接go run 或者 go build 程序会有问题，要修改代码中import有关data包的路径，换成自己机器上的路径

## 依赖
1. golang
2. postgresql
3. 还有一个第三方库包：github.com/lib/pq，用go get即可

## 运行
1. 首先把源代码下载下来：
```shell
git clone https://github.com/liu-jianhao/chitchat.git
```
2. 根据实际情况修改代码中的data包的路径，还有data目录下的data.go中连接数据库的参数
3. 创建数据表（以用户名为chitchat为例，在data目录下执行）：
```shell
psql -U chitchat -f setup.sql
```
4. 返回main.go所在的目录：
```shell
go build
```
若无错误则会得到一个名为chitchat的可执行文件，执行以下语句即可运行：
```shell
./chitchat
```
5. 之后便可在浏览器中输入：
```shell
http://127.0.0.1:8080/
```
即可进入论坛：
![](https://github.com/liu-jianhao/chitchat/tree/master/img/chitchat.png)

## 建议阅读顺序
1. main.go 总体了解整个项目的框架
2. data/user.go 和 data/thread.go 了解chitchat项目中的四个主要结构包含什么成员，特别要注意session结构
3. route_....go 这三个文件主要功能就是处理HTTP请求并响应，session结构发挥了重要作用，特别是其中的UUID成员
4. utils.go 里面有一些在整个项目中用到的函数，在之前遇到其中的函数可以先看
5. 最后是data包中的文件，主要是那四个结构的方法，里面有数据库的操作
