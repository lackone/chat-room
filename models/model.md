### 用户表

```json
{
  "account": "账号",
  "password": "密码",
  "nickname": "昵称",
  "sex": 1,
  "email": "邮箱",
  "avatar": "头像",
  "created_at": 1,
  "updated_at": 1
}
```

### 消息表

```json
{
  "user_identity": "用户的唯一标识",
  "room_identity": "房间的唯一标识",
  "data": "发送的数据",
  "created_at": 1,
  "updated_at": 1
}
```

### 房间表

```json
{
  "number": "房间号",
  "name": "房间名称",
  "info": "房间简介",
  "user_identity": "房间创建者的唯一标识",
  "created_at": 1,
  "updated_at": 1
}
```

### 用户房间关联表

```json
{
  "user_identity": "用户的唯一标识",
  "room_identity": "房间的唯一标识",
  "room_type": "房间类型",
  "created_at": 1,
  "updated_at": 1
}
```
