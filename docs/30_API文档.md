# Agora-BBS MVP API 接口契约文档

## 通用约定

### 1. 响应数据结构

所有 API 统一采用 JSON 格式返回，顶层结构定义如下：

```json
{
  "code": 200,        // 200 表示成功，非 200 表示业务异常错误码
  "message": "success", // 状态描述或错误提示信息
  "data": {}          // 业务数据实体（数组或对象），无数据时为 null
}
```

### 2. 身份认证

需要认证的接口请在请求头（HTTP Header）中带上 JWT Token：

```http
Authorization: Bearer <your_jwt_token>
```

## 1. 用户认证模块 (Auth)

### 1.1 用户注册

- **路由**：`POST /api/v1/auth/register`
    
      
    
- **权限**：公开
    
      
    

**请求参数 (Body)**：

  

JSON

```
{
  "username": "sunny",
  "password": "password123",
  "email": "sunny@example.com"
}
```

**响应示例**：

  

JSON

```
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "username": "sunny",
    "email": "sunny@example.com",
    "trust_score": 100,
    "created_at": "2026-09-07T21:00:00Z"
  }
}
```

### 1.2 用户登录

- **路由**：`POST /api/v1/auth/login`
    
      
    
- **权限**：公开
    
      
    

**请求参数 (Body)**：

  

JSON

```
{
  "username": "sunny",
  "password": "password123"
}
```

**响应示例**：

  

JSON

```
{
  "code": 200,
  "message": "success",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6...",
    "user": {
      "id": 1,
      "username": "sunny",
      "trust_score": 100
    }
  }
}
```

## 2. 主题帖模块 (Topics)

### 2.1 获取主题帖列表（分页）

- **路由**：`GET /api/v1/topics`
    
      
    
- **权限**：公开
    
      
    
- **Query 参数**：`page`（默认 1）、`page_size`（默认 20）、`category_id`（可选筛选）
    
      
    

**响应示例**：

  

JSON

```
{
  "code": 200,
  "message": "success",
  "data": {
    "list": [
      {
        "id": 1,
        "category_id": 1,
        "author_id": 1,
        "author_name": "sunny",
        "title": "关于论坛架构的讨论",
        "status": "published",
        "view_count": 42,
        "reply_count": 5,
        "created_at": "2026-09-07T20:30:00Z"
      }
    ],
    "pagination": {
      "page": 1,
      "page_size": 20,
      "total": 1
    }
  }
}
```

### 2.2 获取主题帖详情

- **路由**：`GET /api/v1/topics/:id`
    
      
    
- **权限**：公开
    
      
    

**响应示例**：

  

JSON

```
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "category_id": 1,
    "author_id": 1,
    "author_name": "sunny",
    "title": "关于论坛架构的讨论",
    "content": "这是一个基于 Go + Next.js + Temporal 的论坛系统...",
    "structured_content": "{}",
    "status": "published",
    "view_count": 43,
    "reply_count": 5,
    "created_at": "2026-09-07T20:30:00Z",
    "updated_at": "2026-09-07T20:30:00Z"
  }
}
```

### 2.3 发布主题帖

- **路由**：`POST /api/v1/topics`
    
      
    
- **权限**：需登录认证
    
      
    

**请求参数 (Body)**：

  

JSON

```
{
  "category_id": 1,
  "title": "关于论坛架构的讨论",
  "content": "这是一个基于 Go + Next.js + Temporal 的论坛系统..."
}
```

**响应示例**：

  

JSON

```
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 1,
    "category_id": 1,
    "author_id": 1,
    "title": "关于论坛架构的讨论",
    "content": "这是一个基于 Go + Next.js + Temporal 的论坛系统...",
    "status": "published",
    "created_at": "2026-09-07T20:30:00Z"
  }
}
```

## 3. 回复/楼层模块 (Posts)

### 3.1 发表回复（支持楼中楼）

- **路由**：`POST /api/v1/topics/:id/posts`
    
      
    
- **权限**：需登录认证
    
      
    

**请求参数 (Body)**：

  

JSON

```
{
  "parent_id": null, // 若回复楼顶帖填 null；若回复某个楼层/楼中楼填对应的 post_id
  "content": "非常赞同这个架构设计！"
}
```

**响应示例**：

JSON

```
{
  "code": 200,
  "message": "success",
  "data": {
    "id": 101,
    "topic_id": 1,
    "author_id": 2,
    "parent_id": null,
    "content": "非常赞同这个架构设计！",
    "post_type": "reply",
    "status": "published",
    "created_at": "2026-09-07T21:15:00Z"
  }
}
```