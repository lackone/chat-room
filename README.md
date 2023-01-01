## 简介

使用 gin+mongo+redis+websocket 开发的简单聊天室后端

## 环境的安装

### 安装扩展

```
go get -u github.com/gin-gonic/gin
go get -u github.com/gorilla/websocket
go get -u github.com/golang-jwt/jwt/v4
go get -u go.mongodb.org/mongo-driver/mongo
go get -u github.com/jordan-wright/email
go get -u github.com/go-redis/redis/v9
```

### 安装mongodb

```
docker run -d --name mongodb -e MONGO_INITDB_ROOT_USERNAME=admin -e MONGO_INITDB_ROOT_PASSWORD=123456 -p 27017:27017 mongo
```

### 安装redis

```
docker run -d --name redis -p 6379:6379 redis --requirepass "123456" --appendonly yes
```